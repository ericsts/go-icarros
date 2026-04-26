# go-icarros

> [🇧🇷 Leia em Português](README.md)

REST API in Go for managing users, cars, and auctions — with JWT authentication, message queues, email notifications, structured event logging, and real-time WebSocket updates.

## Technologies

| Technology | Purpose |
|---|---|
| **Go** + [Gin](https://github.com/gin-gonic/gin) | HTTP framework |
| **PostgreSQL** | database |
| **JWT** (`golang-jwt/jwt`) | stateless authentication |
| **bcrypt** | password hashing |
| **RabbitMQ** (`amqp091-go`) | async event queue |
| **MailHog** + `net/smtp` | email delivery (dev) |
| **WebSocket** (`gorilla/websocket`) | real-time bidding |
| **Docker / Docker Compose** | full dev environment |
| **go-sqlmock** | database mock for tests |

## Architecture

```
cmd/
  main.go                  # entrypoint — wires all dependencies
  seed/main.go             # initial data seeder
internal/
  db/                      # PostgreSQL connection
  models/                  # structs: User, Car, Auction, Bid, EventLog
  repository/              # SQL queries (User, Car, Auction, Bid, Log)
  service/                 # business logic + interfaces
  handler/                 # HTTP handlers + route registration
  middleware/              # AuthMiddleware, AdminMiddleware
  jobs/                    # scheduled tasks (auction auto-close)
  ws/                      # WebSocket hub with per-auction rooms
```

The layers follow: `Handler → Service → Repository → Database`.
Services and handlers depend only on interfaces — no test requires a real database.

### Event flow

```
POST /cars
  └─► AuctionService.CreateForCar   (creates auction in DB)
  └─► Queue.Publish("car.created")  (RabbitMQ)
        └─► Consumer → EmailService  (email to admin)

POST /auctions/:id/bids
  └─► AuctionService.PlaceBid       (validates and saves bid)
  └─► Hub.Broadcast                 (WebSocket → connected clients)

background job (every 1 min)
  └─► AuctionRepository.FindExpired
  └─► AuctionRepository.UpdateStatus("closed")
  └─► Queue.Publish("auction.closed")
        └─► Consumer → EmailService  (email to winner)
```

## How to run

### With Docker Compose (recommended)

```bash
cp .env.example .env   # adjust variables if needed
docker-compose up --build
```

Services available after startup:

| Service | URL |
|---|---|
| API | `http://localhost:8080` |
| MailHog (email UI) | `http://localhost:8025` |
| RabbitMQ (queue panel) | `http://localhost:15672` (guest/guest) |

### Seed — create initial user

```bash
docker-compose exec app go run ./cmd/seed
```

Creates an admin user with email `ericsts@gmail.com` and password `admin123` (configurable in the seeder).

### Without Docker (local database)

```bash
psql -U postgres -c "CREATE DATABASE goapi;"
psql -U postgres -d goapi -f init.sql
cp .env.example .env   # fill in DB_HOST=localhost, etc.
go run ./cmd
```

## Tests

```bash
go test ./...
```

All unit tests — no real database, no RabbitMQ. Each package uses interfaces and dedicated mocks.

---

## Environment variables

Create a `.env` from `.env.example`:

| Variable | Description | Docker default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `db` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | database user | `postgres` |
| `DB_PASSWORD` | database password | `postgres` |
| `DB_NAME` | database name | `goapi` |
| `RABBITMQ_URL` | AMQP connection URL | `amqp://guest:guest@rabbitmq:5672/` |
| `SMTP_HOST` | SMTP host | `mailhog` |
| `SMTP_PORT` | SMTP port | `1025` |
| `SMTP_FROM` | email sender address | `noreply@icarros.com` |
| `ADMIN_EMAIL` | admin notification recipient | `admin@icarros.com` |

---

## Authentication flow

```
POST /login  →  { email, password }
             ←  { token: "eyJ..." }

Protected routes:
  Authorization: Bearer <token>

WebSocket (token via query string):
  GET /ws/auctions/:id?token=<token>
```

The JWT carries `user_id` and `role`. The admin middleware requires `role == "admin"`.

---

## Endpoints

### Public

#### `POST /login`

**Body:**
```json
{ "email": "admin@icarros.com", "password": "password123" }
```
**Response 200:**
```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
```
**Response 401** — invalid credentials.

---

### Users — admin only

> `Authorization: Bearer <admin-token>`

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
**Response 201:** user object (without the `password` field).

#### `GET /users` — list users
**Response 200:** array of users.

#### `GET /users/:id` — get by ID
**Response 200:** user | **404** not found.

#### `PUT /users/:id` — update
**Body:** same format as POST. **Response 200:** updated user.

#### `DELETE /users/:id` — delete
**Response 204** — no body.

---

### Cars — authenticated user

> `Authorization: Bearer <token>`

#### `POST /cars` — register a car + create auction

The `user_id` comes from the token. Registering a car automatically creates an auction.

**Body:**
```json
{
  "marca": "Volkswagen",
  "modelo": "Gol",
  "ano": 2020,
  "valor": 45000.00,
  "auction_ends_at": "2025-12-31T23:59:59Z",
  "min_bid": 40000.00
}
```

> Field names follow Portuguese conventions: `marca` (brand), `modelo` (model), `ano` (year), `valor` (price).

**Response 201:** car object.

> After registration, the admin receives an email notification via RabbitMQ.

#### `GET /cars` — list all cars
**Response 200:** array of cars.

#### `GET /cars/my` — my cars
Returns only the authenticated user's cars.

#### `GET /cars/:id` — get by ID
**Response 200:** car | **404** not found.

#### `PUT /cars/:id` — update
**Body:** car fields. **Response 200:** updated car.

#### `DELETE /cars/:id` — delete
**Response 204.** Admin receives an email notification.

---

### Auctions — authenticated user

> `Authorization: Bearer <token>`

#### `GET /auctions` — list auctions

**Response 200:**
```json
[
  {
    "id": 1,
    "car_id": 3,
    "ends_at": "2025-12-31T23:59:59Z",
    "status": "open",
    "min_bid": 40000,
    "current_bid": 43000,
    "total_bids": 5,
    "created_at": "2025-06-01T10:00:00Z"
  }
]
```

> `current_bid` and `total_bids` are computed in real time via JOIN with the `bids` table.

#### `GET /auctions/:id` — auction by ID
**Response 200:** auction | **404** not found.

#### `POST /auctions/:id/bids` — place a bid

The `user_id` comes from the token. The bid is validated: the auction must be open, not expired, and the amount must exceed the current highest bid (or `min_bid` if no bids exist yet).

**Body:**
```json
{ "amount": 46000.00 }
```
**Response 201:**
```json
{
  "id": 7,
  "auction_id": 1,
  "user_id": 2,
  "amount": 46000,
  "created_at": "2025-06-01T14:32:00Z"
}
```
**Response 400** — auction closed, expired, or bid too low.

> The bid is broadcast via WebSocket to all clients connected to the auction room.

#### `GET /auctions/:id/bids` — bids for an auction
**Response 200:** array of bids ordered highest to lowest.

---

### Logs — admin only

> `Authorization: Bearer <admin-token>`

#### `GET /logs` — list system events

Accepts optional query string filters:

| Parameter | Type | Description |
|---|---|---|
| `level` | string | `info`, `warn`, or `error` |
| `event` | string | filter by event name (ILIKE) |
| `limit` | int | max records (default: 100) |

**Examples:**
```
GET /logs
GET /logs?level=error
GET /logs?event=auction&limit=20
GET /logs?level=info&event=car.created&limit=50
```

**Response 200:**
```json
[
  {
    "id": 1,
    "level": "info",
    "event": "car.created",
    "message": "carro cadastrado",
    "metadata": { "car_id": 3, "user_id": 2 },
    "created_at": "2025-06-01T10:00:00Z"
  }
]
```

---

### WebSocket — real-time bidding

#### `GET /ws/auctions/:id?token=<jwt>`

Connects to the auction's live channel. Authentication is done via the `token` query param (same JWT used for HTTP routes).

Once connected, the client automatically receives each new bid as JSON:

```json
{
  "id": 7,
  "auction_id": 1,
  "user_id": 2,
  "amount": 46000,
  "created_at": "2025-06-01T14:32:00Z"
}
```

Clients are grouped by room (`auction_id`) and only receive messages — data sent by the client is discarded.

---

## Route summary

| Method | Route | Auth | Description |
|--------|-------|:---:|-------------|
| POST | `/login` | public | authenticate and get token |
| POST | `/users` | admin | create user |
| GET | `/users` | admin | list users |
| GET | `/users/:id` | admin | get user |
| PUT | `/users/:id` | admin | update user |
| DELETE | `/users/:id` | admin | delete user |
| GET | `/logs` | admin | list system logs |
| POST | `/cars` | authenticated | register car + create auction |
| GET | `/cars` | authenticated | list all cars |
| GET | `/cars/my` | authenticated | my cars |
| GET | `/cars/:id` | authenticated | get car |
| PUT | `/cars/:id` | authenticated | update car |
| DELETE | `/cars/:id` | authenticated | delete car |
| GET | `/auctions` | authenticated | list auctions |
| GET | `/auctions/:id` | authenticated | get auction |
| POST | `/auctions/:id/bids` | authenticated | place a bid |
| GET | `/auctions/:id/bids` | authenticated | bids for an auction |
| GET | `/ws/auctions/:id?token=` | token via query | WebSocket — live bids |

## Common errors

| Status | Meaning |
|--------|---------|
| 400 | invalid body, malformed parameter, or insufficient bid |
| 401 | missing, invalid, or expired token |
| 403 | valid token but no admin permission |
| 404 | resource not found |
| 500 | internal server error |
