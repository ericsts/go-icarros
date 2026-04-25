# go-icarros

> [🇺🇸 Read in English](README.en.md)

API REST em Go para gerenciamento de usuários e carros, com autenticação JWT e controle de acesso por perfil.

## Tecnologias

- **Go** com [Gin](https://github.com/gin-gonic/gin) — framework HTTP
- **PostgreSQL** — banco de dados
- **JWT** (`golang-jwt/jwt`) — autenticação stateless
- **bcrypt** — hash de senhas
- **Docker / Docker Compose** — ambiente de desenvolvimento
- **go-sqlmock** — mock de banco nos testes

## Arquitetura

```
cmd/
  main.go                  # entrypoint
internal/
  db/                      # conexão com o banco
  models/                  # structs User e Car
  repository/              # queries SQL (UserRepository, CarRepository)
  service/                 # regras de negócio (UserService, CarService, JWT)
  handler/                 # handlers HTTP + registro de rotas
  middleware/              # AuthMiddleware, AdminMiddleware
```

A separação segue o fluxo: `Handler → Service → Repository → Banco`.  
Services e handlers dependem de interfaces, o que permite testes unitários sem banco real.

## Como rodar

### Com Docker Compose (recomendado)

```bash
docker-compose up --build
```

A API sobe em `http://localhost:8080` e o PostgreSQL na porta `5433`.

### Sem Docker (banco local)

Crie o banco e rode o `init.sql`:

```bash
psql -U postgres -c "CREATE DATABASE goapi;"
psql -U postgres -d goapi -f init.sql
```

Ajuste a connection string em `internal/db/db.go` se necessário, depois:

```bash
go run ./cmd
```

## Testes

```bash
go test ./...
```

Os testes são totalmente unitários (sem banco real) usando interfaces e mocks.

---

## Fluxo de autenticação

```
POST /login  →  recebe { email, password }
             →  retorna { token: "eyJ..." }

Requisições protegidas precisam do header:
  Authorization: Bearer <token>
```

O token JWT contém `user_id` e `role`. O middleware de admin verifica se `role == "admin"`.

---

## Endpoints

### Público

#### `POST /login`

Autentica um usuário e retorna o token JWT.

**Body:**
```json
{
  "email": "admin@icarros.com",
  "password": "senha123"
}
```

**Resposta 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Resposta 401** — credenciais inválidas:
```json
{ "error": "credenciais inválidas" }
```

---

### Usuários — apenas admin

Todas as rotas abaixo exigem:
```
Authorization: Bearer <token-de-admin>
```

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

> Valores válidos para `role`: `"user"` ou `"admin"`

**Resposta 201:**
```json
{
  "id": 1,
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "role": "user"
}
```

> A senha nunca é retornada nas respostas.

**Resposta 403** — token não é de admin.

---

#### `GET /users` — listar usuários

**Resposta 200:**
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

#### `GET /users/:id` — buscar usuário por ID

**Resposta 200:**
```json
{
  "id": 1,
  "name": "Eric Santos",
  "email": "eric@icarros.com",
  "role": "user"
}
```

**Resposta 404** — usuário não encontrado.

---

#### `PUT /users/:id` — atualizar usuário

**Body:**
```json
{
  "name": "Eric Luiz",
  "email": "eric.luiz@icarros.com",
  "role": "admin"
}
```

**Resposta 200:**
```json
{
  "id": 1,
  "name": "Eric Luiz",
  "email": "eric.luiz@icarros.com",
  "role": "admin"
}
```

---

#### `DELETE /users/:id` — remover usuário

**Resposta 204** — sem corpo.

---

### Carros — usuário autenticado

Todas as rotas abaixo exigem:
```
Authorization: Bearer <token>
```

Qualquer usuário logado (independente de `role`) pode acessar.

---

#### `POST /cars` — cadastrar carro

O `user_id` é extraído automaticamente do token — não precisa ser enviado no body.

**Body:**
```json
{
  "marca": "Volkswagen",
  "modelo": "Gol",
  "ano": 2020,
  "valor": 45000.00
}
```

**Resposta 201:**
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

#### `GET /cars` — listar todos os carros

**Resposta 200:**
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

#### `GET /cars/my` — listar meus carros

Retorna apenas os carros do usuário autenticado.

**Resposta 200:**
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

#### `GET /cars/:id` — buscar carro por ID

**Resposta 200:**
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

**Resposta 404** — carro não encontrado.

---

#### `PUT /cars/:id` — atualizar carro

**Body:**
```json
{
  "marca": "Volkswagen",
  "modelo": "Polo",
  "ano": 2022,
  "valor": 85000.00
}
```

**Resposta 200:**
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

#### `DELETE /cars/:id` — remover carro

**Resposta 204** — sem corpo.

---

## Resumo das rotas

| Método | Rota | Autenticação | Descrição |
|--------|------|--------------|-----------|
| POST | `/login` | pública | autenticar e obter token |
| POST | `/users` | admin | criar usuário |
| GET | `/users` | admin | listar usuários |
| GET | `/users/:id` | admin | buscar usuário |
| PUT | `/users/:id` | admin | atualizar usuário |
| DELETE | `/users/:id` | admin | remover usuário |
| POST | `/cars` | autenticado | cadastrar carro |
| GET | `/cars` | autenticado | listar todos os carros |
| GET | `/cars/my` | autenticado | meus carros |
| GET | `/cars/:id` | autenticado | buscar carro |
| PUT | `/cars/:id` | autenticado | atualizar carro |
| DELETE | `/cars/:id` | autenticado | remover carro |

## Erros comuns

| Status | Significado |
|--------|-------------|
| 400 | body inválido ou parâmetro mal formatado |
| 401 | token ausente, inválido ou expirado |
| 403 | token válido mas sem permissão de admin |
| 404 | recurso não encontrado |
| 500 | erro interno do servidor |
