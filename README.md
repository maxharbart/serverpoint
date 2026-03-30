# ServerPoint

A multi-Linux server monitoring web application. Monitor CPU, RAM, disk usage, and running processes across your servers via SSH from a single dashboard.

## Features

- Real-time server monitoring via SSH (every 30 seconds)
- Dashboard with server tiles showing CPU/RAM sparklines, disk usage, uptime
- Detailed server view with time-series charts (CPU, RAM), disk usage bar, and top processes table
- AES-GCM encrypted credential storage
- JWT authentication with first-run admin setup
- Dark mode UI

## Quick Start

```bash
# Clone and configure
cp .env.example .env
# Edit .env with your own JWT_SECRET and AES_KEY

# Start all services
docker-compose up --build

# Open http://localhost:3000
# Create your admin account on first run
# Add servers to start monitoring
```

## Architecture

| Service   | Tech                      | Port |
|-----------|---------------------------|------|
| Backend   | Go + Gin + GORM           | 8080 |
| Frontend  | Next.js + TailwindCSS     | 3000 |
| Database  | PostgreSQL 16             | 5432 |

## API Endpoints

| Method | Path                          | Description                    |
|--------|-------------------------------|--------------------------------|
| GET    | /api/auth/check               | Check if setup is required     |
| POST   | /api/auth/register            | Create admin (first-run only)  |
| POST   | /api/auth/login               | Login                          |
| GET    | /api/servers                  | List servers                   |
| POST   | /api/servers                  | Add server                     |
| GET    | /api/servers/:id              | Get server details             |
| DELETE | /api/servers/:id              | Delete server                  |
| GET    | /api/servers/:id/metrics      | Get metrics (query: duration)  |
| GET    | /api/servers/:id/metrics/latest | Get latest metric            |
| GET    | /api/servers/:id/processes    | Get top processes              |

## Environment Variables

| Variable             | Default                            | Description                  |
|----------------------|------------------------------------|------------------------------|
| DB_USER              | serverpoint                        | PostgreSQL username          |
| DB_PASSWORD          | serverpoint                        | PostgreSQL password          |
| DB_NAME              | serverpoint                        | PostgreSQL database name     |
| JWT_SECRET           | change-me-in-production            | JWT signing secret           |
| AES_KEY              | 0123456789abcdef0123456789abcdef   | 32-char hex AES-256 key      |
| NEXT_PUBLIC_API_URL  | http://localhost:8080/api          | Backend API URL for frontend |
