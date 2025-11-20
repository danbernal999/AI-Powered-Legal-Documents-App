package mcp

import (
	"encoding/json"
	"fmt"
)

type MCPProvider struct {
	templates map[string]*DocumentTemplate
	aiClient  AIClient
}

type DocumentTemplate struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	Type      string                   `json:"type"`
	Prompt    string                   `json:"prompt"`
	Variables []TemplateVariable       `json:"variables"`
	Sections  []DocumentSection        `json:"sections"`
}

type TemplateVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
}

type DocumentSection struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type AIClient interface {
	GenerateText(prompt string) (string, error)
	AnalyzeText(text string) (AnalysisResult, error)
}

type AnalysisResult struct {
	Risks      []string `json:"risks"`
	Suggestions []string `json:"suggestions"`
	Score      float64  `json:"score"`
}

type GenerateDocumentInput struct {
	TemplateID string                 `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
}

type GenerateDocumentOutput struct {
	Content string `json:"content"`
	Errors  []string `json:"errors,omitempty"`
}

func NewMCPProvider(aiClient AIClient) *MCPProvider {
	provider := &MCPProvider{
		templates: make(map[string]*DocumentTemplate),
		aiClient:  aiClient,
	}
	provider.initializeTemplates()
	return provider
}

func (p *MCPProvider) initializeTemplates() {
	p.templates["nda"] = &DocumentTemplate{
		ID:   "nda",
		Name: "Non-Disclosure Agreement",
		Type: "nda",
		Prompt: `Generate a comprehensive Non-Disclosure Agreement (NDA) with the following details:
- Disclosing Party: {{disclosing_party}}
- Receiving Party: {{receiving_party}}
- Confidential Information: {{confidential_info}}
- Duration: {{duration}}
- Jurisdiction: {{jurisdiction}}

The NDA should include:
1. Definitions of Confidential Information
2. Obligations of the Receiving Party
3. Permitted Disclosures
4. Term and Termination
5. Return or Destruction of Information
6. Governing Law and Jurisdiction
7. Severability
8. Entire Agreement

Format as a professional legal document.`,
		Variables: []TemplateVariable{
			{Name: "disclosing_party", Description: "Name of the party disclosing information", Type: "string", Required: true},
			{Name: "receiving_party", Description: "Name of the party receiving information", Type: "string", Required: true},
			{Name: "confidential_info", Description: "Description of confidential information", Type: "string", Required: true},
			{Name: "duration", Description: "Duration of confidentiality (e.g., 2 years)", Type: "string", Required: true},
			{Name: "jurisdiction", Description: "Governing jurisdiction", Type: "string", Required: true},
		},
	}

	p.templates["employment"] = &DocumentTemplate{
		ID:   "employment",
		Name: "Employment Contract",
		Type: "employment",
		Prompt: `Generate a comprehensive Employment Contract with the following details:
- Employer: {{employer}}
- Employee: {{employee}}
- Position: {{position}}
- Salary: {{salary}}
- Start Date: {{start_date}}
- Employment Type: {{employment_type}}
- Jurisdiction: {{jurisdiction}}

The contract should include:
1. Position and Responsibilities
2. Compensation and Benefits
3. Term of Employment
4. Working Hours
5. Confidentiality
6. Termination Conditions
7. Non-Competition Clause
8. Governing Law

Format as a professional legal document.`,
		Variables: []TemplateVariable{
			{Name: "employer", Description: "Name of the employer", Type: "string", Required: true},
			{Name: "employee", Description: "Name of the employee", Type: "string", Required: true},
			{Name: "position", Description: "Job position", Type: "string", Required: true},
			{Name: "salary", Description: "Annual salary", Type: "string", Required: true},
			{Name: "start_date", Description: "Employment start date", Type: "string", Required: true},
			{Name: "employment_type", Description: "Full-time or Part-time", Type: "string", Required: true},
			{Name: "jurisdiction", Description: "Governing jurisdiction", Type: "string", Required: true},
		},
	}

	p.templates["rental"] = &DocumentTemplate{
		ID:   "rental",
		Name: "Lease Agreement",
		Type: "rental",
		Prompt: `Generate a comprehensive Lease Agreement with the following details:
- Landlord: {{landlord}}
- Tenant: {{tenant}}
- Property Address: {{property_address}}
- Rent Amount: {{rent_amount}}
- Lease Duration: {{lease_duration}}
- Deposit: {{deposit}}
- Jurisdiction: {{jurisdiction}}

The lease should include:
1. Property Description
2. Rent and Payment Terms
3. Security Deposit
4. Lease Term
5. Use of Property
6. Maintenance and Repairs
7. Utilities
8. Termination and Renewal
9. Governing Law

Format as a professional legal document.`,
		Variables: []TemplateVariable{
			{Name: "landlord", Description: "Name of the landlord", Type: "string", Required: true},
			{Name: "tenant", Description: "Name of the tenant", Type: "string", Required: true},
			{Name: "property_address", Description: "Complete property address", Type: "string", Required: true},
			{Name: "rent_amount", Description: "Monthly rent amount", Type: "string", Required: true},
			{Name: "lease_duration", Description: "Duration of lease (e.g., 12 months)", Type: "string", Required: true},
			{Name: "deposit", Description: "Security deposit amount", Type: "string", Required: true},
			{Name: "jurisdiction", Description: "Governing jurisdiction", Type: "string", Required: true},
		},
	}

	p.templates["freelance"] = &DocumentTemplate{
		ID:   "freelance",
		Name: "Freelance Agreement",
		Type: "freelance",
		Prompt: `Generate a comprehensive Freelance Agreement with the following details:
- Client: {{client}}
- Freelancer: {{freelancer}}
- Scope of Work: {{scope_of_work}}
- Payment Amount: {{payment_amount}}
- Payment Terms: {{payment_terms}}
- Project Duration: {{project_duration}}
- Jurisdiction: {{jurisdiction}}

The agreement should include:
1. Services Description
2. Compensation and Payment Terms
3. Intellectual Property Rights
4. Confidentiality
5. Termination Conditions
6. Liability Limitations
7. Independent Contractor Status
8. Governing Law

Format as a professional legal document.`,
		Variables: []TemplateVariable{
			{Name: "client", Description: "Name of the client", Type: "string", Required: true},
			{Name: "freelancer", Description: "Name of the freelancer", Type: "string", Required: true},
			{Name: "scope_of_work", Description: "Description of services", Type: "string", Required: true},
			{Name: "payment_amount", Description: "Total payment amount", Type: "string", Required: true},
			{Name: "payment_terms", Description: "Payment schedule", Type: "string", Required: true},
			{Name: "project_duration", Description: "Project duration", Type: "string", Required: true},
			{Name: "jurisdiction", Description: "Governing jurisdiction", Type: "string", Required: true},
		},
	}
}

func (p *MCPProvider) GenerateDocument(input GenerateDocumentInput) (GenerateDocumentOutput, error) {
	template, exists := p.templates[input.TemplateID]
	if !exists {
		return GenerateDocumentOutput{
			Errors: []string{fmt.Sprintf("Template %s not found", input.TemplateID)},
		}, fmt.Errorf("template not found")
	}

	prompt := p.buildPrompt(template, input.Variables)

	content, err := p.aiClient.GenerateText(prompt)
	if err != nil {
		return GenerateDocumentOutput{
			Errors: []string{err.Error()},
		}, err
	}

	return GenerateDocumentOutput{
		Content: content,
	}, nil
}

func (p *MCPProvider) buildPrompt(template *DocumentTemplate, variables map[string]interface{}) string {
	prompt := template.Prompt
	for _, v := range template.Variables {
		placeholder := "{{" + v.Name + "}}"
		if val, exists := variables[v.Name]; exists {
			prompt = fmt.Sprintf("%s -> %v", prompt, val)
		} else if v.Required {
			prompt = fmt.Sprintf("%s -> [MISSING REQUIRED: %s]", prompt, v.Name)
		}
	}
	return prompt
}

func (p *MCPProvider) GetTemplate(templateID string) (*DocumentTemplate, error) {
	template, exists := p.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found")
	}
	return template, nil
}

func (p *MCPProvider) GetAllTemplates() []DocumentTemplate {
	templates := make([]DocumentTemplate, 0, len(p.templates))
	for _, t := range p.templates {
		templates = append(templates, *t)
	}
	return templates
}

func (p *MCPProvider) AnalyzeDocument(content string) (AnalysisResult, error) {
	return p.aiClient.AnalyzeText(content)
}

func (p *MCPProvider) GenerateDocumentJSON(input GenerateDocumentInput) (string, error) {
	output, err := p.GenerateDocument(input)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
