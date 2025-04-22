# 📈 Quotes API

Quotes API is a currency quote service that supports asynchronous updates and retrieval of exchange rates for currency pairs.

## 🚀 Features

- Asynchronous quote update by currency pair
- Get quote by ID
- Get the latest quote by currency pair
- Integration with external exchange rates API
- Swagger documentation available at `/swagger/index.html`

## 🛠️ Tech Stack

- Go (Golang)
- Gin (REST API)
- PostgreSQL
- Cron (By robfig)
- Swagger / Swaggo
- Docker / Docker Compose

---

## 🧾 Project Structure

```
plata/
├── build/plata/Dockerfile.go # Dockerfile configuration
├── cmd/plata/main.go         # Application entry point
├── config/config.yaml        # yaml configuration files
├── deployment/
│   └── docker-compose/       # Docker Compose configuration
│       └── docker-compose.yaml
├── docs/                     # Generated Swagger documentation
├── internal/
│   ├── app/
│   │    ├── postgres/        # PostgreSQL client
│   │    └── cron/            # Cron job for updating quotes
│   ├── transport/api/        # HTTP Rest API
│   ├── common/               # Common utilities
│   │    └── log/             # Logger
│   ├── services/quote/       # Business logic for quotes
│   ├── domain/quote/         # data models
│   ├── repository/quote/     # PostgreSQL repository
│   ├── config/               # Load configuration
│   └── clients/exchange/     # External exchange rates API client
├── build/plata/Dockerfile    # Dockerfile for building the application
├── test                      # Unit-tests
├── migrations/               # Database SQL migrations
├── go.mod                    # Go dependencies file
└── README.md                 # Project documentation
```

---

## ⚙️ Build & Run

### 1. Clone the repository

```bash
git clone https://github.com/Russelgon/plata.git
cd plata
```
IMPORTANT: Make sure there is an own API KEY in config.yaml for the external exchange rates API.
### 2. Build and run with Docker

```bash
docker-compose -f deployment/docker-compose/docker-compose.yaml up --build
```

After starting:

- API available at: `http://localhost:8080/api/v1`
- Swagger docs: `http://localhost:8080/swagger/index.html`
- PostgreSQL database will be initialized with the required schema

### 3. Generate Swagger docs manually (if needed)

```bash
swag init --generalInfo cmd/plata/main.go --output docs --parseDependency --parseInternal
```

---

## 📬 API Examples

### Update a quote

```http
POST /api/v1/quotes/update
Headers:
  Content-Type: application/json
  Idempotency-Key: unique-key-123
Body:
{
  "currency": "EUR/USD"
}
```

### Get quote by ID

```http
GET /api/v1/quotes/{id}
```

### Get the latest quote

```http
GET /api/v1/quotes/latest?currency=EUR/USD
```

---

## 🧪 Testing

To run unit tests:

```bash
go test ./...
```

---

## 📝 License

This project is licensed under the MIT License.

