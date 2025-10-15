# Currency Exchange Rate Microservice

A Go microservice that fetches currency exchange rates from the Bank of Latvia RSS feed and provides API endpoints to access current and historical exchange rate data.

## Features

- **Fetches exchange rates** from [Bank of Latvia RSS feed](https://www.bank.lv/vk/ecb_rss.xml)
- **Concurrent processing** using Go routines for 10 preselected currencies
- **RESTful API** with two endpoints for accessing exchange rate data
- **MariaDB database** for persistent storage
- **Extensive logging** with structured JSON logs
- **Docker containerization** for easy deployment
- **Database migrations** using Goose

## API Endpoints

### 1. Get Latest Exchange Rates
```
GET /api/v1/rates/latest
```

Returns the most recent exchange rates for all currencies from the database.

**Response:**
```json
{
  "rates": [
    {
      "currency_code": "SEK",
      "rate": 11.0195,
      "effective_date": "2025-10-15"
    }
  ]
}
```

### 2. Get Currency History
```
GET /api/v1/rates/history/{currency}
```

Returns the complete historical exchange rates for a specific currency.

**Parameters:**
- `currency` - 3-letter currency code (e.g., USD, GBP, SEK)

**Response:**
```json
{
  "currency_code": "SEK",
  "rates": [
    {
      "rate": 11.0195,
      "update_date": "2025-10-15"
    },
    {
      "rate": 11.038,
      "update_date": "2025-10-14"
    },
    {
      "rate": 11.013,
      "update_date": "2025-10-13"
    },
    {
      "rate": 11.008,
      "update_date": "2025-10-10"
    }
  ]
}
```

## Console Commands

### 1. Fetch Command
```bash
./app fetch
```

Fetches current exchange rates from the Bank of Latvia RSS feed and saves them to the database. Uses Go routines to process 10 preselected currencies concurrently.

**Supported Currencies:**
- USD (US Dollar)
- GBP (British Pound)
- JPY (Japanese Yen)
- CHF (Swiss Franc)
- CAD (Canadian Dollar)
- AUD (Australian Dollar)
- SEK (Swedish Krona)
- THB (Thai Baht)
- PLN (Polish Zloty)
- SEK (Swedish Krona)

### 2. Server Command
```bash
./app server
```

Starts the HTTP server and makes the API endpoints accessible to users.

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)
- Make (optional, for using Makefile commands)
- Port 8080 is available for use 

## Quick Start with Docker

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Weeping-Willow/tet.git
   cd tet
   ```

2. **Start the services:**
   ```bash
   docker-compose up
   ```

   This will:
   - Start MariaDB database
   - Run database migrations
   - Fetch initial exchange rates
   - Start the API server on port 8080

3. **Test the API:**
   ```bash
   # Get latest rates
   curl http://localhost:8080/api/v1/rates/latest
   
   # Get USD history
   curl http://localhost:8080/api/v1/rates/history/USD
   ```

### Setup

**Set environment variables:**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_NAME=localhost
   export DB_USER=user
   export DB_PASSWORD=password
   export HTTP_PORT=8080
   ```

   can also use .env file


### Running Commands

```bash
# Build the application
go build -o app cmd/app/main.go

# Fetch exchange rates
./app fetch

# Start the server
./app server
```

## Make Commands

The project includes a Makefile with useful commands:

```bash
# Install dependencies
make deps

# Start all services
make run

# Start only the fetch service, DB is expected to be already running
make run-fetch
```

## Project Structure

```
├── cmd/app/main.go          # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers and middleware
│   ├── app/                 # Application configuration and setup
│   ├── config/              # Configuration management
│   ├── rates/               # Exchange rate business logic
│   ├── storage/             # Database operations and migrations
│   └── utils/               # Utility functions and logging
├── .docker/Dockerfile       # Docker build configuration
├── docker-compose.yml       # Docker services definition
├── go.mod                   # Go module dependencies
└── Makefile                # Build automation
```