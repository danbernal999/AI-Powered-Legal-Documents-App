# KiraDoc API Examples

## Authentication

### 1. Register New User

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123",
    "name": "Daniel Bernal"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "Daniel Bernal"
}
```

---

### 2. Login

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "Daniel Bernal"
  }
}
```

**Save the token** - you'll need it for all subsequent requests.

---

## Templates

### 3. Get All Templates

**Request:**
```bash
curl http://localhost:8080/api/v1/templates \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```json
[
  {
    "id": "nda",
    "name": "Non-Disclosure Agreement",
    "description": "Protect confidential information with a professional NDA",
    "type": "nda",
    "variables": [
      {
        "name": "disclosing_party",
        "description": "Name of the party disclosing information",
        "type": "string",
        "required": true
      },
      {
        "name": "confidentiality_period",
        "description": "Duration of confidentiality in years",
        "type": "number",
        "required": true
      }
      // ... more variables
    ]
  },
  {
    "id": "employment",
    "name": "Employment Contract",
    "description": "Create a professional employment agreement",
    "type": "employment",
    "variables": [...]
  },
  {
    "id": "rental",
    "name": "Lease Agreement",
    "description": "Professional property rental agreement",
    "type": "rental",
    "variables": [...]
  },
  {
    "id": "freelance",
    "name": "Freelance Agreement",
    "description": "Independent contractor service agreement",
    "type": "freelance",
    "variables": [...]
  }
]
```

---

### 4. Get Specific Template

**Request:**
```bash
curl http://localhost:8080/api/v1/templates/nda \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```json
{
  "id": "nda",
  "name": "Non-Disclosure Agreement",
  "description": "Protect confidential information with a professional NDA",
  "type": "nda",
  "variables": [
    {
      "name": "disclosing_party",
      "description": "Name of the party disclosing information",
      "type": "string",
      "required": true
    },
    {
      "name": "receiving_party",
      "description": "Name of the party receiving information",
      "type": "string",
      "required": true
    },
    // ... more variables
  ]
}
```

---

## Document Generation

### 5. Generate Document from Template

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "nda",
    "title": "Acme Corp & Tech Startup NDA",
    "variables": {
      "disclosing_party": "Acme Corporation",
      "receiving_party": "TechStartup Inc",
      "execution_date": "2024-11-17",
      "disclosing_entity_type": "Corporation",
      "receiving_entity_type": "Limited Liability Company",
      "confidential_info": "AI Technology, Trade Secrets, Source Code",
      "duration": "3",
      "jurisdiction": "California"
    }
  }'
```

**Response:**
```json
{
  "id": "doc-123456",
  "content": "NON-DISCLOSURE AGREEMENT\n\nThis Non-Disclosure Agreement (\"Agreement\") is entered into as of 2024-11-17, between Acme Corporation, a Corporation (\"Disclosing Party\"), and TechStartup Inc, a Limited Liability Company (\"Receiving Party\")...\n\n1. DEFINITIONS\n   1.1 \"Confidential Information\" means..."
}
```

---

## Document Management

### 6. Save Document

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "nda",
    "title": "Acme Corp & Tech Startup NDA",
    "type": "nda",
    "content": "NON-DISCLOSURE AGREEMENT\n\nThis Non-Disclosure Agreement...",
    "variables": {
      "disclosing_party": "Acme Corporation",
      "receiving_party": "TechStartup Inc",
      "execution_date": "2024-11-17",
      "confidential_info": "AI Technology, Trade Secrets, Source Code",
      "duration": "3",
      "jurisdiction": "California"
    }
  }'
```

**Response:**
```json
{
  "id": "doc-550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user-550e8400-e29b-41d4-a716-446655440001",
  "template_id": "nda",
  "title": "Acme Corp & Tech Startup NDA",
  "type": "nda",
  "content": "NON-DISCLOSURE AGREEMENT...",
  "variables": {...},
  "status": "draft",
  "version": 1,
  "created_at": "2024-11-17T10:30:00Z",
  "updated_at": "2024-11-17T10:30:00Z"
}
```

---

### 7. Get All Documents

**Request:**
```bash
curl http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```json
[
  {
    "id": "doc-550e8400-e29b-41d4-a716-446655440000",
    "user_id": "user-550e8400-e29b-41d4-a716-446655440001",
    "template_id": "nda",
    "title": "Acme Corp & Tech Startup NDA",
    "type": "nda",
    "content": "NON-DISCLOSURE AGREEMENT...",
    "variables": {...},
    "status": "draft",
    "version": 1,
    "created_at": "2024-11-17T10:30:00Z",
    "updated_at": "2024-11-17T10:30:00Z"
  },
  {
    "id": "doc-550e8400-e29b-41d4-a716-446655440002",
    "user_id": "user-550e8400-e29b-41d4-a716-446655440001",
    "template_id": "employment",
    "title": "Employment Contract - Jane Smith",
    "type": "employment",
    "content": "EMPLOYMENT AGREEMENT...",
    "variables": {...},
    "status": "draft",
    "version": 1,
    "created_at": "2024-11-17T11:15:00Z",
    "updated_at": "2024-11-17T11:15:00Z"
  }
]
```

---

### 8. Get Specific Document

**Request:**
```bash
curl http://localhost:8080/api/v1/documents/doc-550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```json
{
  "id": "doc-550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user-550e8400-e29b-41d4-a716-446655440001",
  "template_id": "nda",
  "title": "Acme Corp & Tech Startup NDA",
  "type": "nda",
  "content": "NON-DISCLOSURE AGREEMENT...",
  "variables": {...},
  "status": "draft",
  "version": 1,
  "created_at": "2024-11-17T10:30:00Z",
  "updated_at": "2024-11-17T10:30:00Z"
}
```

---

### 9. Update Document

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/documents/doc-550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated NDA Title",
    "content": "MODIFIED CONTENT...",
    "status": "finalized",
    "version": 2,
    "variables": {...}
  }'
```

**Response:**
```json
{
  "id": "doc-550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated NDA Title",
  "content": "MODIFIED CONTENT...",
  "status": "finalized",
  "version": 2,
  "updated_at": "2024-11-17T12:00:00Z"
}
```

---

### 10. Delete Document

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/documents/doc-550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

**Response:**
```
HTTP 204 No Content
```

---

## Complete Workflow Example

### Step 1: Register & Login

```bash
# Register
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "workflow@example.com",
    "password": "TestPass123",
    "name": "Workflow User"
  }')

# Login & extract token
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "workflow@example.com",
    "password": "TestPass123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"
```

### Step 2: Fetch Templates

```bash
curl -s http://localhost:8080/api/v1/templates \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Step 3: Generate NDA

```bash
curl -s -X POST http://localhost:8080/api/v1/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "nda",
    "title": "My First NDA",
    "variables": {
      "disclosing_party": "My Company",
      "receiving_party": "Partner Company",
      "execution_date": "2024-11-17",
      "disclosing_entity_type": "Corporation",
      "receiving_entity_type": "Corporation",
      "confidential_info": "Business Plans and Technology",
      "duration": "3",
      "jurisdiction": "New York"
    }
  }' | jq
```

### Step 4: Save Document

```bash
curl -s -X POST http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "nda",
    "title": "My First NDA",
    "type": "nda",
    "content": "NON-DISCLOSURE AGREEMENT...",
    "variables": {
      "disclosing_party": "My Company",
      "receiving_party": "Partner Company"
    }
  }' | jq
```

---

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Invalid token"
}
```

### 404 Not Found
```json
{
  "error": "Template not found"
}
```

### 400 Bad Request
```json
{
  "error": "Invalid request"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to generate document"
}
```

---

## Tips

1. **Always include the Authorization header** for protected endpoints
2. **Use `jq`** to format JSON responses: `curl ... | jq`
3. **Save tokens** in environment variables for easier testing
4. **Check timestamps** in responses to verify operations
5. **Use document IDs** returned from create operations for updates/deletes
