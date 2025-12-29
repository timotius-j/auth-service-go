# Auth Service (Go)

Auth Service is a lightweight authentication API built with Go, Gin, and PostgreSQL. It supports registration, login, and a password reset flow using OTP + reset tokens. The codebase follows a clean architecture style (handler -> service -> repository) and includes local migrations and seed data.

## Features
- User registration and login with bcrypt password hashing
- JWT access token generation
- Password reset flow: request OTP, verify OTP, reset password with token
- PostgreSQL persistence with migrations and seed data (local only)
- Structured logging and centralized error handling

## Tech Stack
- Language: Go
- Web framework: Gin
- Database: PostgreSQL
- Driver: pgx
- Logging: Zap

## Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── dto/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── route/
│   │   └── service/
│   ├── apperror/
│   ├── common/
│   ├── config/
│   ├── dto/
│   ├── middleware/
│   └── util/
├── migrations/
├── pkg/
├── .env.example
├── go.mod
└── README.md
```

## API
Base path: `/api/v1`

- GET `/auth/` - Get user data by email (expects JSON body)
- POST `/auth/register` - Register a new user
- POST `/auth/login` - Login and receive a JWT
- POST `/auth/forgot-password` - Request a reset OTP
- POST `/auth/verify-reset-password` - Verify OTP and receive a reset token
- POST `/auth/reset-password` - Reset password using reset token

Response envelope:

```json
{
  "success": true,
  "data": {}
}
```

Error envelope:

```json
{
  "success": false,
  "error": {
    "message": "Invalid Request"
  }
}
```

Example (register):

```json
POST /api/v1/auth/register
{
  "username": "user1",
  "email": "user1@example.com",
  "password": "password123"
}
```

```json
{
  "success": true,
  "data": {
    "message": "Registration successful. Please log in to continue."
  }
}
```

## Configuration
Copy `.env.example` to `.env` and fill in values.

Required environment variables:
- `JWT_SECRET` - Signing key for JWT tokens
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE` - PostgreSQL connection
- `APP_ENV` - Use `local` to enable auto migrations and seed data

Password reset configuration:
- `RESET_OTP_TTL_MINUTES`
- `RESET_TOKEN_TTL_MINUTES`
- `RESET_OTP_MAX_ATTEMPTS`
- `RESET_OTP_RESEND_COOLDOWN_SECONDS`
- `RESET_OTP_RESEND_LIMIT_PER_HOUR`
- `RESET_OTP_HMAC_SECRET`

Note: OTP delivery uses a dummy email sender that logs the OTP to stdout.

## Running Locally
1. `cp .env.example .env`
2. Update `.env` with your database and reset settings
3. Ensure PostgreSQL is running
4. Start the server:

```bash
go run ./cmd/server
```

The server listens on `:8080`.

## Migrations and Seed Data
Migrations and seed data run automatically on startup when `APP_ENV=local`. The scripts live in `migrations/` and include:
- `000_drop_tables.sql`
- `001_schema.sql`
- `002_seeder.sql`
