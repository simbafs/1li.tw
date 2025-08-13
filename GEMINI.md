## Project Overview

This is a full-stack URL shortener application named `1li.tw`.

**Backend:**
- **Language:** Go
- **Framework:** Gin
- **Database:** SQLite
- **ORM/Query Builder:** sqlc
- **Authentication:** JWT
- **Features:**
    - User registration and login
    - URL shortening
    - Custom paths for shortened URLs
    - Analytics for shortened URLs
    - Telegram bot integration for managing URLs

**Frontend:**
- **Framework:** Astro
- **Styling:** Tailwind CSS with daisyUI
- **Location:** `web/` directory

The application is structured following a clean architecture pattern, with distinct layers for domain, application, infrastructure, and presentation.

## Building and Running

### Backend (Go)

To run the backend server:

```bash
# From the root directory
go run main.go
```

The server will start on port 8080 by default. You can change this in your `.env` file or by setting the `SERVER_PORT` environment variable.

### Frontend (Astro)

To run the frontend development server:

```bash
# From the web/ directory
cd web
pnpm install
pnpm run dev
```

To build the frontend for production:

```bash
# From the web/ directory
cd web
pnpm run build
```

## Development Conventions

### Database

- The database schema is defined in `sql/schema.sql`.
- SQL queries are located in `sql/queries/`.
- `sqlc` is used to generate type-safe Go code from the SQL schema and queries. To regenerate the code, run:

```bash
sqlc generate
```

### API

The backend exposes a RESTful API under the `/api` prefix. Authentication is handled via JWTs.

### Code Style

- **Go:** Follows standard Go conventions.
- **Frontend:** Uses Prettier for code formatting. To format the frontend code, run:

```bash
# From the web/ directory
cd web
pnpm run format
```
