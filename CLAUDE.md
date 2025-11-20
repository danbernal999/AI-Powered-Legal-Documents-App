# KiraDoc - Developer Quick Reference

## Project Overview

**KiraDoc** is an AI-powered legal document generator that helps users create customized contracts (NDAs, employment agreements, etc.) using intelligent templates and AI.

**Stack:**
- Backend: Go with Gorilla Mux, PostgreSQL, JWT, MCP
- Frontend: Next.js 14, React 18, Zustand, Tailwind CSS
- Deployment: Docker + Docker Compose

---

## Project Structure

```
AI-Powered-Legal-Documents-App/
├── backend/
│   ├── cmd/main.go              # Entry point
│   ├── pkg/
│   │   ├── handlers/            # HTTP handlers
│   │   ├── services/            # Business logic
│   │   ├── repository/          # Database operations + seeds
│   │   ├── middleware/          # JWT middleware
│   │   ├── mcp/                 # MCP provider + AI integration
│   │   └── models/              # Data models
│   ├── go.mod                   # Go dependencies
│   ├── Dockerfile               # Docker configuration
│   └── .env                     # Environment variables
├── frontend/
│   ├── app/                     # Next.js app directory
│   │   ├── page.tsx            # Home page
│   │   ├── auth/                # Login/Register
│   │   ├── dashboard/           # User dashboard
│   │   └── documents/           # Document pages
│   ├── lib/
│   │   ├── api.ts              # Axios API client
│   │   └── store.ts            # Zustand stores
│   ├── package.json
│   ├── tsconfig.json
│   └── .env.local
├── docker-compose.yml           # Docker services
├── README.md                    # Full documentation
├── SETUP.md                     # Setup instructions
└── .env                         # Backend env

```

---

## Key Features Implemented

✅ User authentication (JWT)
✅ PostgreSQL database with migrations
✅ 4 legal templates: NDA, Employment, Rental, Freelance
✅ API endpoints for templates, documents, generation
✅ MCP provider for AI integration
✅ OpenAI/Groq client implementations
✅ Next.js frontend with dynamic forms
✅ Document CRUD operations
✅ Template seeding on startup
✅ Responsive UI with Tailwind CSS

---

## Running the Project

### Option 1: Docker Compose (Recommended)
```bash
docker-compose up
```

Then in another terminal:
```bash
cd frontend
npm install
npm run dev
```

### Option 2: Manual Setup
```bash
# Backend
cd backend
go mod download
go run cmd/main.go

# Frontend (new terminal)
cd frontend
npm install
npm run dev
```

**Access:**
- Frontend: http://localhost:3000
- API: http://localhost:8080/api/v1
- Database: localhost:5432

---

## Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://user:password@localhost:5432/kiradoc
PORT=8080
JWT_SECRET=your-secret-key
OPENAI_API_KEY=your-key (optional)
GROQ_API_KEY=your-key (optional)
```

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

---

## API Endpoints

### Auth
- `POST /auth/register` - Create user
- `POST /auth/login` - Login user

### Templates
- `GET /templates` - Get all templates
- `GET /templates/{id}` - Get specific template

### Documents
- `GET /documents` - List user documents
- `POST /documents` - Create document
- `GET /documents/{id}` - Get document
- `PUT /documents/{id}` - Update document
- `DELETE /documents/{id}` - Delete document

### Generation
- `POST /generate` - Generate from template

**All endpoints except auth require JWT token in Authorization header**

---

## Database Schema

```sql
-- Users
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255),
  name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Templates
CREATE TABLE templates (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255),
  description TEXT,
  type VARCHAR(50),
  content TEXT,
  variables TEXT (JSON),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Documents
CREATE TABLE documents (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) REFERENCES users,
  template_id VARCHAR(36) REFERENCES templates,
  title VARCHAR(255),
  type VARCHAR(50),
  content TEXT,
  variables TEXT (JSON),
  status VARCHAR(50),
  version INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

---

## Available Templates

1. **NDA** (id: `nda`)
   - Variables: disclosing_party, receiving_party, confidentiality_period, etc.

2. **Employment** (id: `employment`)
   - Variables: company_name, employee_name, job_title, salary, etc.

3. **Rental Lease** (id: `rental`)
   - Variables: landlord_name, tenant_name, monthly_rent, etc.

4. **Freelance Agreement** (id: `freelance`)
   - Variables: client_name, contractor_name, project_fee, etc.

---

## Common Development Tasks

### Add a new template
1. Edit `backend/pkg/repository/seeds.go`
2. Add template to the `templates` slice
3. Restart backend (migrations won't re-run if tables exist)

### Fix frontend errors
- Check browser console (F12)
- Check .env.local for correct API URL
- Ensure backend is running

### Fix backend errors
- Check PostgreSQL connection
- Verify DATABASE_URL in .env
- Run `go mod tidy` if import errors

### Add AI integration
- Set `OPENAI_API_KEY` or `GROQ_API_KEY` in .env
- Already implemented in `pkg/mcp/ai_client.go`
- Will auto-generate documents using AI

---

## Next Steps for Improvement

### High Priority
- [ ] Integrate real AI model calls (OpenAI/Groq)
- [ ] Implement PDF export
- [ ] Add document analysis for legal risks
- [ ] Improve error handling and validation

### Medium Priority
- [ ] Mobile app (Flutter)
- [ ] Electronic signatures
- [ ] Document collaboration
- [ ] User settings/preferences
- [ ] Search and filtering

### Low Priority
- [ ] Multi-language support
- [ ] Advanced analytics
- [ ] Audit logs
- [ ] Admin dashboard

---

## Testing Quick Commands

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass","name":"Test"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass"}'

# Get templates
curl http://localhost:8080/api/v1/templates \
  -H "Authorization: Bearer <TOKEN>"

# Generate document
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id":"nda",
    "title":"My NDA",
    "variables":{...}
  }'
```

---

## Known Issues & Workarounds

| Issue | Solution |
|-------|----------|
| Port 8080 in use | Change PORT in .env |
| DB connection error | Verify PostgreSQL running, check DATABASE_URL |
| Module not found (Go) | Run `go mod download` in backend |
| Frontend blank | Check NEXT_PUBLIC_API_URL in .env.local |
| JWT errors | Verify JWT_SECRET matches between backend and middleware |

---

## Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [JWT Introduction](https://jwt.io/introduction)

---

## Last Updated

- Backend: ✅ Complete with MVP features
- Frontend: ✅ Complete with MVP features
- Documentation: ✅ Complete

**Status:** Ready for testing and AI integration
