# Senior Full-Stack Developer Challenge: Task Management Platform

## Overview
Build a collaborative task management platform with real-time features, advanced filtering, and team collaboration capabilities. This challenge tests your ability to architect scalable solutions, implement complex business logic, and create intuitive user experiences.

## Core Requirements

### Backend — Go (Preferred)

**Preferred Stack: Go with Gin or Echo**

We prefer Go for its performance, clean concurrency model, and production-ready binaries. However, you may use **Node.js/Express** or **Python/FastAPI** if you can justify your choice in your SOLUTION.md.

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin, Manager, Member)
  - Password reset functionality
  - Rate limiting on auth endpoints

- **Database Design**
  - Design schema for users, teams, projects, tasks, and comments
  - Implement proper relationships and constraints
  - Add database migrations
  - Include seed data for testing

- **API Endpoints**
  - RESTful API with proper HTTP status codes
  - CRUD operations for all entities
  - Advanced filtering (by status, assignee, date range, tags) and pagination
  - File upload for task attachments
  - Bulk operations (assign multiple tasks, update statuses)

- **Real-time Features**
  - WebSocket implementation for live updates
  - Real-time notifications for task assignments

### Frontend (React/Vue/Angular or your choice)
- **User Interface**
  - Responsive design (mobile-first approach)
  - Dark/light theme toggle
  - Drag-and-drop task management (Kanban board)
  - Advanced search and filtering
  - Data visualization (charts for task completion, team performance)

- **State Management**
  - Implement proper state management
  - Optimistic updates for better UX
  - Cache management for improved performance

- **Advanced Features**
  - Infinite scrolling for task lists
  - Keyboard shortcuts for power users
  - Export functionality (CSV data minimum, PDF reports bonus)

---

## MVP vs Full Implementation

### MVP Tier (3 days) — Must Have
These are the **minimum requirements** for a passing submission:

| Area | Requirements |
|------|-------------|
| **Auth** | JWT login/register, basic RBAC (Admin, Member) |
| **API** | CRUD for users, projects, tasks with pagination |
| **Frontend** | Task list view, create/edit tasks, basic filtering |
| **Database** | Migrations, seed data, proper relationships |
| **DevOps** | Docker Compose that works with `docker-compose up` |
| **Testing** | Unit tests for critical business logic |

### Full Implementation (5-7 days) — Should Have
Build on MVP with these additions:

| Area | Requirements |
|------|-------------|
| **Auth** | Password reset, rate limiting, Manager role |
| **API** | File uploads, bulk operations, comments |
| **Real-time** | WebSocket notifications for task assignments |
| **Frontend** | Kanban board (drag-and-drop), dark/light theme, charts |
| **Performance** | Redis caching, query optimization |
| **Testing** | Integration tests, 70%+ coverage |

### Excellence Tier — Nice to Have
Bonus features that demonstrate senior-level thinking:

- Live collaboration indicators (who's viewing what)
- Offline capability with sync
- Time tracking for tasks
- Email notifications / digest
- Third-party integrations (GitHub/Slack webhooks)
- AI-powered task categorization
- Audit logging
- Load testing results

---

## Technical Challenges

### 1. Performance & Scalability
- Implement database query optimization
- Add caching layer (Redis recommended)
- Lazy loading for large datasets
- Image optimization for uploaded files
- API response compression

### 2. Security
- Input validation and sanitization
- SQL injection prevention
- XSS protection
- CSRF tokens
- Secure file upload handling
- API versioning strategy

### 3. Testing
- Unit tests for business logic
- Integration tests for API endpoints
- E2E tests for critical user flows (bonus)
- Target: **70% coverage** on critical paths

### 4. DevOps & Deployment
- **Docker containerization** (Required)
  - Multi-stage Docker builds for optimization
  - Docker Compose for local development
  - Separate containers for frontend, backend, database
- Environment configuration management
- CI/CD pipeline setup (GitHub Actions)
- Error monitoring and logging

---

## Evaluation Criteria

### Code Quality (25%)
- Clean, readable, and maintainable code
- Proper error handling and edge cases
- Consistent coding standards
- Documentation and comments

### Architecture (25%)
- Scalable and modular design
- Proper separation of concerns
- Database design efficiency
- API design best practices

### Functionality (25%)
- All MVP features working correctly
- User experience and interface design
- Real-time features implementation
- Performance optimization

### Technical Excellence (25%)
- Security implementation
- Testing coverage and quality
- DevOps and deployment setup
- Bonus features implementation

---

## Getting Started

### Project Structure Suggestion

```
task-management-platform/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── pkg/
│   ├── tests/
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   ├── public/
│   ├── tests/
│   ├── Dockerfile
│   └── package.json
├── database/
│   ├── migrations/
│   └── seeds/
├── docs/
├── docker-compose.yml
├── docker-compose.prod.yml
└── README.md
```

### Docker Requirements

#### Local Development Setup
Your solution **must** include Docker configuration for easy local setup:

```bash
# Should work out of the box
git clone <your-fork>
cd <your-fork>
docker-compose up
```

#### Required Docker Files
- **docker-compose.yml**: For local development environment
- **docker-compose.prod.yml**: For production deployment (optional)
- **Backend Dockerfile**: Multi-stage build optimized for Go
- **Frontend Dockerfile**: Optimized build for your chosen frontend
- **Environment Variables**: Properly configured for different environments

#### Container Requirements
- **Frontend Container**: Serve built static files or run dev server
- **Backend Container**: API server with health checks
- **Database Container**: PostgreSQL with persistent volumes
- **Redis Container**: For caching and session management
- **Reverse Proxy**: Nginx container for routing (bonus)

All containers should start with a single `docker-compose up` command and be accessible via localhost with clear port documentation.

---

## Sample User Stories
1. As a team member, I want to see all my assigned tasks in a dashboard
2. As a manager, I want to create projects and assign tasks to team members
3. As a user, I want to receive real-time notifications when tasks are updated
4. As an admin, I want to view team performance analytics

---

## How to Submit Your Solution

### 1. Fork and Setup
1. **Fork this repository** to your GitHub account
2. **Clone your fork** locally and start building
3. **Organize your code** following the structure above

### 2. Build Your Solution
- Follow the requirements outlined above
- Commit regularly with clear commit messages
- Focus on code quality and documentation

### 3. Document Your Work
Create a **SOLUTION.md** file containing:
- **Setup Instructions**: How to run locally (Docker required)
- **Technology Stack**: What you used and why (especially if not using Go)
- **Architecture Overview**: High-level design decisions
- **API Documentation**: Endpoints and usage
- **Database Schema**: Structure and relationships
- **Demo Credentials**: Test accounts for different roles
- **Live Demo URL**: Deployed application link
- **Known Issues**: Any limitations or incomplete features
- **Time Spent**: Rough breakdown of hours per area

### 4. Deploy Your Application
Deploy to any free hosting service:
- **Frontend**: Vercel, Netlify, GitHub Pages
- **Backend**: Railway, Render, Fly.io
- **Database**: Supabase, Neon, PlanetScale

### 5. Final Submission
When ready, provide us with:
- **GitHub repository link** (your fork)
- **Live demo URL**
- **Any additional context** or notes

---

## Questions?

### Using GitHub Issues (Recommended)
We've set up issue templates to help you get quick, targeted assistance:

1. **Go to the [Issues tab](../../issues)** in this repository
2. **Click "New issue"**
3. **Choose the appropriate template:**
   - 🔧 **Technical Question** - Implementation and coding issues
   - ❓ **Requirements Clarification** - Unclear requirements or specifications
   - 📤 **Submission Help** - Process, documentation, or deployment questions
   - 🐛 **Bug Report** - Issues with challenge materials

4. **We'll respond within 24-48 hours** during business days

---

## Technology Notes

### Why We Prefer Go
- **Single binary deployment** — trivial Docker images, no runtime dependencies
- **Built-in concurrency** — goroutines make WebSockets and real-time features elegant
- **Strong typing** — catches issues at compile time
- **Performance** — fast execution, low memory footprint
- **Clean architecture** — the language encourages good design patterns

### Acceptable Alternatives
If you choose Node.js or Python, explain in your SOLUTION.md:
- Why you chose it over Go
- How you addressed performance considerations
- Your approach to type safety (TypeScript, type hints)

---

**Remember**: Quality over quantity. A well-implemented MVP with clean code, good tests, and solid documentation beats a feature-complete mess. Show us how you think, not just what you can copy-paste.
