# go-icarros

> [🇺🇸 Read in English](README.en.md)

API REST em Go para gerenciamento de usuários, carros e leilões — com autenticação JWT, fila de mensagens, notificações por e-mail, logs estruturados e WebSocket em tempo real.

## Tecnologias

| Tecnologia | Uso |
|---|---|
| **Go** + [Gin](https://github.com/gin-gonic/gin) | framework HTTP |
| **PostgreSQL** | banco de dados |
| **JWT** (`golang-jwt/jwt`) | autenticação stateless |
| **bcrypt** | hash de senhas |
| **RabbitMQ** (`amqp091-go`) | fila de eventos assíncronos |
| **MailHog** + `net/smtp` | envio de e-mails (dev) |
| **WebSocket** (`gorilla/websocket`) | lances em tempo real |
| **Docker / Docker Compose** | ambiente completo |
| **go-sqlmock** | mock de banco nos testes |

## Arquitetura

```
cmd/
  main.go                  # entrypoint — wiring de todas as dependências
  seed/main.go             # seeder de dados iniciais
internal/
  db/                      # conexão com PostgreSQL
  models/                  # structs: User, Car, Auction, Bid, EventLog
  repository/              # queries SQL (User, Car, Auction, Bid, Log)
  service/                 # regras de negócio + interfaces
  handler/                 # handlers HTTP + registro de rotas
  middleware/              # AuthMiddleware, AdminMiddleware
  jobs/                    # tarefas agendadas (encerramento de leilões)
  ws/                      # hub WebSocket com salas por leilão
```

O fluxo segue: `Handler → Service → Repository → Banco`.
Services e handlers dependem apenas de interfaces — nenhum teste precisa de banco real.

### Fluxo de eventos

```
POST /cars
  └─► AuctionService.CreateForCar   (cria leilão no banco)
  └─► Queue.Publish("car.created")  (RabbitMQ)
        └─► Consumer → EmailService  (e-mail ao admin)

POST /auctions/:id/bids
  └─► AuctionService.PlaceBid       (valida e salva lance)
  └─► Hub.Broadcast                 (WebSocket → clientes conectados)

job (1 min)
  └─► AuctionRepository.FindExpired
  └─► AuctionRepository.UpdateStatus("closed")
  └─► Queue.Publish("auction.closed")
        └─► Consumer → EmailService  (e-mail ao vencedor)
```

## Como rodar

### Com Docker Compose (recomendado)

```bash
cp .env.example .env   # ajuste as variáveis se necessário
docker-compose up --build
```

Serviços disponíveis após o start:

| Serviço | URL |
|---|---|
| API | `http://localhost:8080` |
| MailHog (UI de e-mails) | `http://localhost:8025` |
| RabbitMQ (painel de filas) | `http://localhost:15672` (guest/guest) |

### Seed — criar usuário inicial

```bash
docker-compose exec app go run ./cmd/seed
```

Cria um usuário admin com e-mail `ericsts@gmail.com` e senha `admin123` (configurável no seeder).

### Sem Docker (banco local)

```bash
psql -U postgres -c "CREATE DATABASE goapi;"
psql -U postgres -d goapi -f init.sql
cp .env.example .env   # preencha DB_HOST=localhost, etc.
go run ./cmd
```

## Testes

```bash
go test ./...
```

Todos unitários — sem banco real, sem RabbitMQ. Cada pacote usa interfaces e mocks dedicados.

---

## Variáveis de ambiente

Crie um `.env` a partir de `.env.example`:

| Variável | Descrição | Padrão Docker |
|---|---|---|
| `DB_HOST` | host do PostgreSQL | `db` |
| `DB_PORT` | porta do PostgreSQL | `5432` |
| `DB_USER` | usuário | `postgres` |
| `DB_PASSWORD` | senha | `postgres` |
| `DB_NAME` | nome do banco | `goapi` |
| `RABBITMQ_URL` | URL de conexão AMQP | `amqp://guest:guest@rabbitmq:5672/` |
| `SMTP_HOST` | host SMTP | `mailhog` |
| `SMTP_PORT` | porta SMTP | `1025` |
| `SMTP_FROM` | remetente dos e-mails | `noreply@icarros.com` |
| `ADMIN_EMAIL` | destinatário das notificações de admin | `admin@icarros.com` |

---

## Fluxo de autenticação

```
POST /login  →  { email, password }
             ←  { token: "eyJ..." }

Rotas protegidas:
  Authorization: Bearer <token>

WebSocket (autenticação via query string):
  GET /ws/auctions/:id?token=<token>
```

O token JWT contém `user_id` e `role`. O middleware de admin exige `role == "admin"`.

---

## Endpoints

### Público

#### `POST /login`

**Body:**
```json
{ "email": "admin@icarros.com", "password": "senha123" }
```
**Resposta 200:**
```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
```
**Resposta 401** — credenciais inválidas.

---

### Usuários — apenas admin

> `Authorization: Bearer <token-de-admin>`

#### `POST /users` — criar usuário

**Body:**
```json
{
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "password": "senha123",
  "role": "user"
}
```
**Resposta 201:** objeto do usuário (sem o campo `password`).

#### `GET /users` — listar usuários
**Resposta 200:** array de usuários.

#### `GET /users/:id` — buscar por ID
**Resposta 200:** usuário | **404** não encontrado.

#### `PUT /users/:id` — atualizar
**Body:** mesmo formato do POST. **Resposta 200:** usuário atualizado.

#### `DELETE /users/:id` — remover
**Resposta 204** — sem corpo.

---

### Carros — usuário autenticado

> `Authorization: Bearer <token>`

#### `POST /cars` — cadastrar carro + criar leilão

O `user_id` vem do token. Ao cadastrar, um leilão é criado automaticamente.

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
**Resposta 201:** objeto do carro.

> Após o cadastro, o admin recebe um e-mail de notificação via RabbitMQ.

#### `GET /cars` — listar todos os carros
**Resposta 200:** array de carros.

#### `GET /cars/my` — meus carros
Retorna apenas os carros do usuário autenticado.

#### `GET /cars/:id` — buscar por ID
**Resposta 200:** carro | **404** não encontrado.

#### `PUT /cars/:id` — atualizar
**Body:** campos do carro. **Resposta 200:** carro atualizado.

#### `DELETE /cars/:id` — remover
**Resposta 204.** Admin recebe e-mail de notificação.

---

### Leilões — usuário autenticado

> `Authorization: Bearer <token>`

#### `GET /auctions` — listar leilões

**Resposta 200:**
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

> `current_bid` e `total_bids` são calculados em tempo real via JOIN com a tabela `bids`.

#### `GET /auctions/:id` — leilão por ID
**Resposta 200:** leilão | **404** não encontrado.

#### `POST /auctions/:id/bids` — dar um lance

O `user_id` vem do token. O lance é validado: o leilão deve estar aberto, não expirado, e o valor deve ser maior que o lance atual (ou o `min_bid` se não houver lances).

**Body:**
```json
{ "amount": 46000.00 }
```
**Resposta 201:**
```json
{
  "id": 7,
  "auction_id": 1,
  "user_id": 2,
  "amount": 46000,
  "created_at": "2025-06-01T14:32:00Z"
}
```
**Resposta 400** — leilão encerrado, expirado ou lance insuficiente.

> O lance é transmitido via WebSocket a todos os clientes conectados no leilão.

#### `GET /auctions/:id/bids` — lances de um leilão
**Resposta 200:** array de lances ordenados do maior para o menor.

---

### Logs — apenas admin

> `Authorization: Bearer <token-de-admin>`

#### `GET /logs` — listar eventos do sistema

Aceita filtros via query string:

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `level` | string | `info`, `warn` ou `error` |
| `event` | string | filtra por nome de evento (ILIKE) |
| `limit` | int | máximo de registros (padrão: 100) |

**Exemplos:**
```
GET /logs
GET /logs?level=error
GET /logs?event=auction&limit=20
GET /logs?level=info&event=car.created&limit=50
```

**Resposta 200:**
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

### WebSocket — lances em tempo real

#### `GET /ws/auctions/:id?token=<jwt>`

Conecta ao canal do leilão. A autenticação é feita via query param `token` (mesmo JWT das rotas HTTP).

Após conectar, o cliente recebe automaticamente cada novo lance no formato JSON:

```json
{
  "id": 7,
  "auction_id": 1,
  "user_id": 2,
  "amount": 46000,
  "created_at": "2025-06-01T14:32:00Z"
}
```

Clientes são agrupados por sala (`auction_id`) e apenas recebem mensagens — o envio de dados pelo cliente é ignorado.

---

## Resumo das rotas

| Método | Rota | Autenticação | Descrição |
|--------|------|:---:|-----------|
| POST | `/login` | pública | autenticar e obter token |
| POST | `/users` | admin | criar usuário |
| GET | `/users` | admin | listar usuários |
| GET | `/users/:id` | admin | buscar usuário |
| PUT | `/users/:id` | admin | atualizar usuário |
| DELETE | `/users/:id` | admin | remover usuário |
| GET | `/logs` | admin | listar logs do sistema |
| POST | `/cars` | autenticado | cadastrar carro + criar leilão |
| GET | `/cars` | autenticado | listar todos os carros |
| GET | `/cars/my` | autenticado | meus carros |
| GET | `/cars/:id` | autenticado | buscar carro |
| PUT | `/cars/:id` | autenticado | atualizar carro |
| DELETE | `/cars/:id` | autenticado | remover carro |
| GET | `/auctions` | autenticado | listar leilões |
| GET | `/auctions/:id` | autenticado | buscar leilão |
| POST | `/auctions/:id/bids` | autenticado | dar um lance |
| GET | `/auctions/:id/bids` | autenticado | lances de um leilão |
| GET | `/ws/auctions/:id?token=` | token via query | WebSocket — lances ao vivo |

## Erros comun

| Status | Significado |
|--------|-------------|
| 400 | body inválido, parâmetro mal formatado ou lance insuficiente |
| 401 | token ausente, inválido ou expirado |
| 403 | token válido mas sem permissão de admin |
| 404 | recurso não encontrado |
| 500 | erro interno do servidor |
