# go-icarros

> [🇧🇷 Leia em Português](README.md)

REST API in Go for managing users and cars, with JWT authentication and role-based access control.

## Technologies

- **Go** with [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- **PostgreSQL** — database
- **JWT** (`golang-jwt/jwt`) — stateless authentication
- **bcrypt** — password hashing
- **Docker / Docker Compose** — development environment
- **go-sqlmock** — database mock for tests

## Architecture

```
cmd/
  main.go                  # entrypoint
internal/
  db/                      # database connection
  models/                  # User and Car structs
  repository/              # SQL queries (UserRepository, CarRepository)
  service/                 # business logic (UserService, CarService, JWT)
  handler/                 # HTTP handlers + route registration
  middleware/              # AuthMiddleware, AdminMiddleware
```

The layers follow the flow: `Handler → Service → Repository → Database`.  
Services and handlers depend on interfaces, enabling unit tests without a real database.

## How to run

### With Docker Compose (recommended)

```bash
docker-compose up --build
```

The API starts at `http://localhost:8080` and PostgreSQL on port `5433`.

### Without Docker (local database)

Create the database and run `init.sql`:

```bash
psql -U postgres -c "CREATE DATABASE goapi;"
psql -U postgres -d goapi -f init.sql
```

Adjust the connection string in `internal/db/db.go` if needed, then:

```bash
go run ./cmd
```

## Tests

```bash
go test ./...
```

All tests are fully unit-tested (no real database) using interfaces and mocks.

---

## Authentication flow

```
POST /login  →  receives { email, password }
             →  returns  { token: "eyJ..." }

Protected requests require the header:
  Authorization: Bearer <token>
```

The JWT token carries `user_id` and `role`. The admin middleware checks whether `role == "admin"`.

---

## Endpoints

### Public

#### `POST /login`

Authenticates a user and returns a JWT token.

**Body:**
```json
{
  "email": "admin@icarros.com",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response 401** — invalid credentials:
```json
{ "error": "credenciais inválidas" }
```

---

### Users — admin only

All routes below require:
```
Authorization: Bearer <admin-token>
```

#### `POST /users` — create user

**Body:**
```json
{
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "password": "password123",
  "role": "user"
}
```

> Valid values for `role`: `"user"` or `"admin"`

**Response 201:**
```json
{
  "id": 1,
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "role": "user"
}
```

> Passwords are never returned in responses.

**Response 403** — token does not belong to an admin.

---

#### `GET /users` — list users

**Response 200:**
```json
[
  {
    "id": 1,
    "name": "Eric Santos",
    "email": "eric@icarros.com",
    "role": "user"
  },
  {
    "id": 2,
    "name": "Ana Lima",
    "email": "ana@icarros.com",
    "role": "admin"
  }
]
```

---

#### `GET /users/:id` — get user by ID

**Response 200:**
```json
{
  "id": 1,
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "role": "user"
}
```

**Response 404** — user not found.

---

#### `PUT /users/:id` — update user

**Body:**
```json
{
  "name": "Eric Luiz",
  "email": "eric.luiz@icarros.com",
  "role": "admin"
}
```

**Response 200:**
```json
{
  "id": 1,
  "name": "Eric Luiz",
  "email": "eric.luiz@icarros.com",
  "role": "admin"
}
```

---

#### `DELETE /users/:id` — delete user

**Response 204** — no body.

---

### Cars — authenticated user

All routes below require:
```
Authorization: Bearer <token>
```

Any logged-in user (regardless of `role`) can access these.

---

#### `POST /cars` — register a car

The `user_id` is automatically extracted from the token — it does not need to be sent in the body.

**Body:**
```json
{
  "marca": "Volkswagen",
  "modelo": "Gol",
  "ano": 2020,
  "valor": 45000.00
}
```

> Field names follow Portuguese conventions: `marca` (brand), `modelo` (model), `ano` (year), `valor` (price).

**Response 201:**
```json
{
  "id": 1,
  "user_id": 3,
  "marca": "Volkswagen",
  "modelo": "Gol",
  "ano": 2020,
  "valor": 45000
}
```

---

#### `GET /cars` — list all cars

**Response 200:**
```json
[
  {
    "id": 1,
    "user_id": 3,
    "marca": "Volkswagen",
    "modelo": "Gol",
    "ano": 2020,
    "valor": 45000
  },
  {
    "id": 2,
    "user_id": 1,
    "marca": "Fiat",
    "modelo": "Uno",
    "ano": 2018,
    "valor": 28000
  }
]
```

---

#### `GET /cars/my` — list my cars

Returns only the cars belonging to the authenticated user.

**Response 200:**
```json
[
  {
    "id": 1,
    "user_id": 3,
    "marca": "Volkswagen",
    "modelo": "Gol",
    "ano": 2020,
    "valor": 45000
  }
]
```

---

#### `GET /cars/:id` — get car by ID

**Response 200:**
```json
{
  "id": 1,
  "user_id": 3,
  "marca": "Volkswagen",
  "modelo": "Gol",
  "ano": 2020,
  "valor": 45000
}
```

**Response 404** — car not found.

---

#### `PUT /cars/:id` — update car

**Body:**
```json
{
  "marca": "Volkswagen",
  "modelo": "Polo",
  "ano": 2022,
  "valor": 85000.00
}
```

**Response 200:**
```json
{
  "id": 1,
  "user_id": 0,
  "marca": "Volkswagen",
  "modelo": "Polo",
  "ano": 2022,
  "valor": 85000
}
```

---

#### `DELETE /cars/:id` — delete car

**Response 204** — no body.

---

## Route summary

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| POST | `/login` | public | authenticate and get token |
| POST | `/users` | admin | create user |
| GET | `/users` | admin | list users |
| GET | `/users/:id` | admin | get user |
| PUT | `/users/:id` | admin | update user |
| DELETE | `/users/:id` | admin | delete user |
| POST | `/cars` | authenticated | register car |
| GET | `/cars` | authenticated | list all cars |
| GET | `/cars/my` | authenticated | my cars |
| GET | `/cars/:id` | authenticated | get car |
| PUT | `/cars/:id` | authenticated | update car |
| DELETE | `/cars/:id` | authenticated | delete car |

## Common errors

| Status | Meaning |
|--------|---------|
| 400 | invalid body or malformed parameter |
| 401 | missing, invalid, or expired token |
| 403 | valid token but no admin permission |
| 404 | resource not found |
| 500 | internal server error |
