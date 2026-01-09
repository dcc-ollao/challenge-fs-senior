# SOLUTION.md

## 1. Setup Instructions (Local)

### Requirements
- Docker
- Docker Compose

### Run locally

    docker-compose up --build

This starts:
- PostgreSQL
- Backend API
- Database migrations
- Frontend application

### URLs
- Frontend: http://localhost:5173  
- Backend: http://localhost:8080  

Environment variables are included so the project runs without additional setup.

---

## 2. Technology Stack

### Backend
- Go
- Gin
- PostgreSQL
- SQLx
- JWT authentication
- bcrypt password hashing
- golang-migrate for schema migrations

Go was chosen for simplicity, performance, and explicit control over business logic.

### Frontend
- React
- TypeScript
- Vite
- Axios
- Tailwind CSS

### Infrastructure
- Frontend deployed on Vercel
- Backend deployed on Render
- PostgreSQL hosted on Render

---

## 3. Architecture Overview

The backend follows a layered structure:

    Handlers → Services → Repositories → Database

- Business rules live in the service layer
- Permissions are enforced server-side
- Repositories expose explicit queries
- Frontend mirrors backend rules to prevent invalid actions

---

## 4. API Documentation

### Authentication
- POST /auth/register
- POST /auth/login
- POST /auth/change-password
- GET /auth/me

JWT is sent via the Authorization header.

---

### Users (Admin only)
- GET /users
- PUT /users/:id
- DELETE /users/:id

---

### Projects
- GET /projects
- POST /projects
- PUT /projects/:id
- DELETE /projects/:id

---

### Tasks
- GET /tasks
- POST /tasks
- PUT /tasks/:id
- DELETE /tasks/:id

Rules:
- Admin can modify and assign any task
- Users can only change status of tasks assigned to them

---

### Admin Export
- GET /admin/export

Exports all application data as a ZIP:
- users.csv
- projects.csv
- tasks.csv

---

## 5. Database Schema

### Users
- id
- email
- password_hash
- role
- created_at

### Projects
- id
- name
- owner_id
- created_at

### Tasks
- id
- project_id
- title
- description
- status
- assignee_id
- created_at
- updated_at

---

## 6. Demo Credentials

    Admin
    admin@test.com
    admin123

    User
    user@test.com
    user123

---

## 7. Live Demo

- Frontend: https://challenge-fs-senior-eta.vercel.app  
- Backend: https://challenge-fs-senior.onrender.com  

---

## 8. Known Issues / Limitations

- Task filtering is client-side
- No background jobs
- No refresh tokens
- Minimal UI styling by design

---

## 9. Time Spent

| Area               | Time     |
|--------------------|----------|
| Backend            | ~6 hours |
| Frontend           | ~4 hours |
| Deployment         | ~2 hours |
| Testing & fixes    | ~2 hours |
| **Total**          | ~14 hours |
