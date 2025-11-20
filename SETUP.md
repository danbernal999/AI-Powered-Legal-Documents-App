# KiraDoc Setup Guide

## Quick Start with Docker Compose

The fastest way to get started is using Docker Compose:

```bash
docker-compose up
```

This will:
1. Start PostgreSQL database on port 5432
2. Start the backend API on port 8080
3. Automatically create database schema and seed templates

Then start the frontend in a new terminal:

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at http://localhost:3000

---

## Manual Setup

### Prerequisites
- Go 1.21+ (Install from https://golang.org/dl/)
- Node.js 18+ (Install from https://nodejs.org/)
- PostgreSQL 15+ (Install from https://www.postgresql.org/download/)

### Step 1: Setup PostgreSQL

```bash
# Create database
createdb kiradoc

# Or use psql:
psql -U postgres
CREATE DATABASE kiradoc;
```

### Step 2: Backend Setup

```bash
cd backend

# Download dependencies
go mod download

# Create .env file with your PostgreSQL connection
echo "DATABASE_URL=postgres://postgres:password@localhost:5432/kiradoc" > ../.env
echo "PORT=8080" >> ../.env
echo "JWT_SECRET=your-super-secret-key" >> ../.env
echo "OPENAI_API_KEY=your-openai-key" >> ../.env

# Run the backend
go run cmd/main.go
```

The backend will:
1. Migrate the database schema
2. Seed legal templates
3. Start the API server on port 8080

### Step 3: Frontend Setup

In a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will be available at http://localhost:3000

---

## First Time Usage

1. Open http://localhost:3000 in your browser
2. Click "Sign Up" to create an account
3. After login, click "Create Document"
4. Select "NDA" template as your first document
5. Fill in the required fields:
   - Disclosing Party: "Acme Corp"
   - Receiving Party: "Tech Startup Inc"
   - Confidential Information: "AI technology and trade secrets"
   - Duration: "3"
   - Jurisdiction: "California"
6. Click "Generate Document"
7. Preview the generated NDA
8. Click "Save Document"

---

## Testing the API with cURL

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get Templates

```bash
curl http://localhost:8080/api/v1/templates \
  -H "Authorization: Bearer <TOKEN_FROM_LOGIN>"
```

### Generate Document

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "template_id": "nda",
    "title": "My First NDA",
    "variables": {
      "disclosing_party": "Acme Corp",
      "receiving_party": "Tech Startup",
      "execution_date": "2024-11-17",
      "disclosing_entity_type": "Corporation",
      "receiving_entity_type": "LLC",
      "confidential_info": "AI Technology",
      "duration": "3",
      "jurisdiction": "California"
    }
  }'
```

---

## Troubleshooting

### PostgreSQL Connection Error
- Make sure PostgreSQL is running
- Verify the DATABASE_URL is correct
- Check that the database `kiradoc` exists

### Port Already in Use
- Backend: Change PORT in .env (default 8080)
- Frontend: npm run dev will prompt to use a different port

### Modules Not Found (Go)
- Run `go mod download` in the backend directory
- Delete go.sum and run `go mod tidy`

### Node Modules Issues
- Delete `node_modules` and `package-lock.json`
- Run `npm install` again

---

## Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://user:password@localhost:5432/kiradoc
PORT=8080
JWT_SECRET=your-secret-key-here
OPENAI_API_KEY=sk-...  (optional for OpenAI)
GROQ_API_KEY=...        (optional for Groq)
MCP_ENABLED=true
```

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

---

## Production Deployment

### Backend on Heroku

```bash
cd backend
heroku create kiradoc-backend
git push heroku main
```

### Frontend on Vercel

```bash
cd frontend
vercel
```

---

## Next Steps

After the MVP is working:
1. Add more templates (services agreement, partnership agreement, etc.)
2. Implement document analysis for legal risk detection
3. Add PDF export with formatting
4. Create Flutter mobile app
5. Implement electronic signatures
6. Add collaboration features

---

## Support

For issues:
1. Check the console logs
2. Review the README.md for detailed documentation
3. Check backend logs: `tail -f backend/logs.txt`
4. Check browser console for frontend errors

Happy document generation! 🚀
