# First Go Project

A CRUD admin dashboard built with Go, HTMX, and Supabase (PostgreSQL).

## Tech Stack

- **Backend**: Go + Chi router
- **Frontend**: HTMX + HTML templates
- **Database**: Supabase (PostgreSQL)

## Getting Started

### Prerequisites

- Go 1.24+
- A Supabase project

### Setup

1. Clone the repo

```bash
   git clone https://github.com/pedwoo/first-go-project.git
   cd first-go-project
```

2. Install dependencies

```bash
   go mod tidy
```

3. Copy the example env file and fill in your values

```bash
   cp .env.example .env
```

4. Run the server

```bash
   go run main.go
```

5. In a separate terminal

```bash
   sass --watch ./static/scss/main.scss ./static/css/main.css
```

5. Open your browser at `http://localhost:8080`

## Project Structure

```
├── main.go           # Entry point
├── db/               # Database connection and queries
├── handlers/         # HTTP route handlers
└── templates/        # HTML templates
    ├── layout/       # Base layout
    ├── pages/        # Full pages
    └── partials/     # HTMX fragments
```
