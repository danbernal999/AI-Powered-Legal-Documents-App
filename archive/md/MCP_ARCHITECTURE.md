# KiraDoc MCP (Model Context Protocol) Architecture

## Overview

The Model Context Protocol (MCP) in KiraDoc is the bridge between user input and AI-powered document generation. It manages templates, variables, and calls to AI models (OpenAI/Groq).

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                        │
│  User fills form → API Request                               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│              Backend Handlers (Go)                           │
│  /generate endpoint receives template_id + variables        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│             MCP Provider (pkg/mcp)                           │
│                                                              │
│  1. Load Template by ID                                      │
│  2. Build Prompt with user variables                         │
│  3. Call AI Client (OpenAI/Groq)                             │
│  4. Return generated document content                        │
└────────────────────┬────────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         ▼                       ▼
┌──────────────────┐  ┌──────────────────┐
│  OpenAI Client   │  │  Groq Client     │
│                  │  │                  │
│ gpt-3.5-turbo    │  │ mixtral-8x7b     │
└──────────────────┘  └──────────────────┘
         │                       │
         └───────────┬───────────┘
                     ▼
         ┌───────────────────────┐
         │  AI Model Response    │
         │  (Generated Document) │
         └───────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │  Return to Handler    │
         │  → Save to Database   │
         │  → Send to Frontend   │
         └───────────────────────┘
```

---

## Component Details

### 1. DocumentTemplate Structure

```go
type DocumentTemplate struct {
    ID        string                  // Unique identifier (e.g., "nda")
    Name      string                  // Display name
    Type      string                  // Template type
    Prompt    string                  // AI prompt with placeholders
    Variables []TemplateVariable      // Required user inputs
    Sections  []DocumentSection       // Document sections
}

type TemplateVariable struct {
    Name        string  // Variable name
    Description string  // User-friendly description
    Type        string  // "string", "number", "date"
    Required    bool    // Is this required?
}
```

### 2. MCPProvider Methods

```go
// Main generation method
func (p *MCPProvider) GenerateDocument(input GenerateDocumentInput) (GenerateDocumentOutput, error)

// Template management
func (p *MCPProvider) GetTemplate(templateID string) (*DocumentTemplate, error)
func (p *MCPProvider) GetAllTemplates() []DocumentTemplate

// AI analysis
func (p *MCPProvider) AnalyzeDocument(content string) (AnalysisResult, error)

// Helper methods
func (p *MCPProvider) buildPrompt(template *DocumentTemplate, variables map[string]interface{}) string
```

### 3. AI Client Interface

```go
type AIClient interface {
    GenerateText(prompt string) (string, error)
    AnalyzeText(text string) (AnalysisResult, error)
}
```

---

## Flow: Document Generation

### Step 1: User Submits Form

```json
{
  "template_id": "nda",
  "title": "My NDA",
  "variables": {
    "disclosing_party": "Acme Corp",
    "receiving_party": "Tech Startup",
    "confidentiality_period": "3",
    "governing_jurisdiction": "California"
  }
}
```

### Step 2: Handler Receives Request

```go
func (h *Handlers) GenerateDocumentHandler(w http.ResponseWriter, r *http.Request) {
    var req models.GenerateDocumentRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Request goes to MCP provider
    output, err := h.mcpProvider.GenerateDocument(req)
}
```

### Step 3: MCP Provider Processes

```go
func (p *MCPProvider) GenerateDocument(input GenerateDocumentInput) (GenerateDocumentOutput, error) {
    // 1. Load template from memory/database
    template := p.templates[input.TemplateID]
    
    // 2. Build prompt by replacing placeholders with user variables
    prompt := p.buildPrompt(template, input.Variables)
    // Example result:
    // "Generate a comprehensive Non-Disclosure Agreement with:
    //  - Disclosing Party: Acme Corp
    //  - Receiving Party: Tech Startup
    //  - Confidentiality Period: 3 years
    //  - Jurisdiction: California"
    
    // 3. Call AI model
    content, err := p.aiClient.GenerateText(prompt)
    
    // 4. Return generated content
    return GenerateDocumentOutput{Content: content}, nil
}
```

### Step 4: AI Model Generates Document

The prompt is sent to OpenAI or Groq, which returns a fully formatted legal document:

```
NON-DISCLOSURE AGREEMENT

This Non-Disclosure Agreement ("Agreement") is entered into as of [date], 
between Acme Corp, a [entity type] ("Disclosing Party"), and Tech Startup, 
a [entity type] ("Receiving Party").

WHEREAS, the Disclosing Party possesses certain confidential information...
```

### Step 5: Response Sent Back

```json
{
  "id": "doc-123456",
  "content": "NON-DISCLOSURE AGREEMENT\n\nThis Non-Disclosure Agreement..."
}
```

---

## Template Structure Example: NDA

```go
&DocumentTemplate{
    ID:   "nda",
    Name: "Non-Disclosure Agreement",
    Type: "nda",
    Prompt: `Generate a comprehensive NDA with:
            - Disclosing Party: {{disclosing_party}}
            - Receiving Party: {{receiving_party}}
            - Confidential Information: {{confidential_info}}
            - Duration: {{duration}}
            - Jurisdiction: {{jurisdiction}}
            Include: definitions, obligations, permitted disclosures, etc.`,
    Variables: []TemplateVariable{
        {
            Name:        "disclosing_party",
            Description: "Name of the party disclosing information",
            Type:        "string",
            Required:    true,
        },
        {
            Name:        "duration",
            Description: "Duration of confidentiality in years",
            Type:        "number",
            Required:    true,
        },
        // ... more variables
    },
}
```

---

## OpenAI Client Implementation

```go
type OpenAIClient struct {
    apiKey string
    model  string  // "gpt-3.5-turbo", "gpt-4", etc.
}

func (c *OpenAIClient) GenerateText(prompt string) (string, error) {
    // Build OpenAI API request
    req := OpenAIRequest{
        Model: c.model,
        Messages: []Message{
            {
                Role: "system",
                Content: "You are a legal document expert. Generate professional legal documents.",
            },
            {
                Role: "user",
                Content: prompt,
            },
        },
    }
    
    // Send to OpenAI API
    resp, _ := http.Post("https://api.openai.com/v1/chat/completions", ...)
    
    // Parse response
    var openaiResp OpenAIResponse
    json.Unmarshal(body, &openaiResp)
    
    // Return generated content
    return openaiResp.Choices[0].Message.Content, nil
}
```

---

## Groq Client Implementation

```go
type GroqClient struct {
    apiKey string
    model  string  // "mixtral-8x7b-32768", etc.
}

func (c *GroqClient) GenerateText(prompt string) (string, error) {
    // Similar to OpenAI but uses Groq API endpoint
    // https://api.groq.com/openai/v1/chat/completions
    
    // Groq is faster and cheaper for many use cases
}
```

---

## Mock AI for Development

When `OPENAI_API_KEY` or `GROQ_API_KEY` are not set, the system uses a mock response:

```go
if c.apiKey == "" {
    return "Mock generated document for: " + prompt, nil
}
```

This allows frontend development without API keys configured.

---

## Adding a New Template

### Step 1: Define Template in seeds.go

```go
{
    id:      "service",
    name:    "Service Agreement",
    typeStr: "service",
    prompt: `Generate a comprehensive Service Agreement with:
            - Service Provider: {{provider_name}}
            - Client: {{client_name}}
            - Services: {{services_description}}
            - Fee: {{service_fee}}`,
    variables: `[
        {"name":"provider_name",...},
        {"name":"client_name",...},
        {"name":"services_description",...},
        {"name":"service_fee",...}
    ]`,
}
```

### Step 2: Restart Backend

Templates are seeded on startup if the table is empty. For existing databases, manually insert:

```sql
INSERT INTO templates (id, name, description, type, content, variables)
VALUES ('service', 'Service Agreement', '...', 'service', '...', '[...]');
```

### Step 3: Use in Frontend

The new template will automatically appear in the template selection.

---

## Document Analysis (Future Feature)

```go
func (p *MCPProvider) AnalyzeDocument(content string) (AnalysisResult, error) {
    // Analyze document for risks
    prompt := fmt.Sprintf(`Analyze this legal document for risks:
                          %s
                          
                          Provide JSON with: risks (array), suggestions (array), score (0-100)`, 
                          content)
    
    content, _ := p.aiClient.GenerateText(prompt)
    
    var result AnalysisResult
    json.Unmarshal([]byte(content), &result)
    
    return result, nil
}
```

---

## API Integration Points

### Generate Endpoint

```
POST /api/v1/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "template_id": "nda",
  "variables": { ... },
  "title": "My Document"
}
```

Handler calls MCP provider → returns generated content

### Templates Endpoint

```
GET /api/v1/templates
Authorization: Bearer <token>
```

Returns list of all templates with their variables

---

## Performance Considerations

1. **Template Caching**: Templates are loaded in memory on startup
2. **AI Calls**: Currently synchronous - consider async for production
3. **Prompt Building**: String replacement is O(n) complexity
4. **Error Handling**: If AI fails, return error - no fallback

---

## Security Considerations

1. **API Keys**: Store in environment variables, never in code
2. **Prompt Injection**: Variables are substituted safely
3. **User Data**: Confidential information is stored encrypted
4. **Rate Limiting**: Implement per-user API call limits

---

## Future Enhancements

- [ ] Cache AI responses for identical inputs
- [ ] Support multiple AI providers simultaneously
- [ ] Implement document versioning in MCP
- [ ] Add A/B testing for different prompts
- [ ] Implement streaming responses for large documents
- [ ] Add webhook callbacks for long-running generations
