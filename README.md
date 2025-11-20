# KiraDoc - AI-Powered Legal Document Generator

KiraDoc is an application that enables users to generate customized legal contracts (NDAs, employment agreements, lease agreements, freelance contracts, etc.) using AI-powered templates.

## Architecture Overview

### Backend (Go)
- **Framework**: Gorilla Mux for routing
- **Database**: PostgreSQL with SQL migrations
- **Authentication**: JWT tokens
- **AI Integration**: OpenAI/Groq via MCP (Model Context Protocol)

**Structure**:
- `cmd/main.go` - Application entry point
- `pkg/handlers/` - HTTP request handlers
- `pkg/services/` - Business logic
- `pkg/repository/` - Database operations
- `pkg/middleware/` - JWT authentication
- `pkg/mcp/` - AI provider and template management

### Frontend (Next.js)
- **Framework**: React 18 with Next.js 14
- **State Management**: Zustand
- **Styling**: Tailwind CSS
- **API Communication**: Axios

**Structure**:
- `app/` - Next.js app router pages
- `lib/` - API client and Zustand stores
- `app/auth/` - Login and registration pages
- `app/dashboard/` - User dashboard
- `app/documents/` - Document creation and editing

## Features

### MVP
- ✅ User authentication (registration/login)
- ✅ Document template selection
- ✅ Dynamic form for contract variables
- ✅ AI-powered document generation
- ✅ Document preview and editing
- ✅ Save and manage documents
- ✅ Document export options

### Available Templates
1. **NDA** (Non-Disclosure Agreement)
2. **Employment Contract**
3. **Lease Agreement**
4. **Freelance Agreement**

## Installation & Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp ../.env .env
```

4. Update the `.env` file with your configuration:
```
DATABASE_URL=
PORT=
JWT_SECRET=your-
OPENAI_API_KEY=
GROQ_API_KEY=
```

5. Run the backend:
```bash
go run cmd/main.go
```

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

4. Open http://localhost:3000 in your browser

### Using Docker Compose

For a complete setup with PostgreSQL:

```bash
docker-compose up
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 8080

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Templates
- `GET /api/v1/templates` - Get all templates
- `GET /api/v1/templates/{id}` - Get specific template

### Documents
- `GET /api/v1/documents` - List user's documents
- `POST /api/v1/documents` - Create/save a document
- `GET /api/v1/documents/{id}` - Get document details
- `PUT /api/v1/documents/{id}` - Update document
- `DELETE /api/v1/documents/{id}` - Delete document

### Generation
- `POST /api/v1/generate` - Generate document from template

## MCP Integration

The backend uses MCP (Model Context Protocol) to:
1. **Manage Templates**: Store and retrieve legal document templates
2. **Handle Variables**: Process user input variables
3. **Generate Documents**: Call OpenAI/Groq to create documents
4. **Analyze Content**: Check for risks and suggest improvements

### MCP Provider Structure

```go
type MCPProvider struct {
    templates map[string]*DocumentTemplate
    aiClient  AIClient
}
```

## Database Schema

### Tables
- `users` - User accounts and authentication
- `templates` - Legal document templates
- `documents` - Generated documents and versions

## Security

- JWT token-based authentication
- Password hashing with bcrypt
- CORS headers configured
- Private document access control

## Future Enhancements

- [ ] Electronic signature integration
- [ ] Multi-user collaboration
- [ ] Advanced legal analysis
- [ ] Document version history
- [ ] PDF export with formatting
- [ ] Email notifications
- [ ] Admin dashboard for templates
- [ ] Mobile app (Flutter)

## Testing

To run tests:

### Backend
```bash
cd backend
go test ./...
```

### Frontend
```bash
cd frontend
npm test
```

## Deployment

### Backend (Docker)
```bash
cd backend
docker build -t kiradoc-backend .
docker run -p 8080:8080 kiradoc-backend
```

### Frontend (Vercel/Netlify)
```bash
cd frontend
npm run build
npm start
```

## Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://user:password@localhost/kiradoc
PORT=8080
JWT_SECRET=your-secret-key
OPENAI_API_KEY=your-openai-key
GROQ_API_KEY=your-groq-key
MCP_ENABLED=true
```

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: [Create an issue]
- Email: support@kiradoc.com

---

**Built with ❤️ by the KiraDoc team**
