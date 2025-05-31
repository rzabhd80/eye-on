# Eye-On 👁️

**Exchange and vendor agnostic transaction manager**

A sophisticated Go-based trading and exchange management system that provides unified access to multiple cryptocurrency exchanges with robust order management, user authentication, and real-time data processing.

## 🚀 Features

- **Multi-Exchange Support**: Unified interface for multiple cryptocurrency exchanges (Bitpin, Nobitex)
- **Order Management**: Comprehensive order tracking, history, and snapshot capabilities
- **User Authentication**: JWT-based authentication with secure credential management
- **Real-time Data**: Redis-powered caching and real-time data processing
- **RESTful API**: Clean HTTP API with middleware for authentication and validation
- **Database Migrations**: Structured database schema management
- **Docker Ready**: Complete containerization with Docker Compose

## 🏗️ Architecture

The project follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Business logic and entities
- **API Layer**: HTTP handlers and routing
- **Infrastructure Layer**: Database, Redis, and external service connections
- **Application Layer**: Use cases and application services

## 📁 Project Structure

```
eye-on/
├── .env.dockercompose.example
├── go.mod
├── go.sum
├── README.md
│
├── cmd/
│   ├── api.go
│   └── app.go
│
├── infra/
│   ├── docker-compose.yml
│   └── Dockerfile
│
├── internal/
│   ├── database/
│   │   └── models/
│   │       ├── baseModel.go
│   │       ├── exchange.go
│   │       ├── exchangeCredentials.go
│   │       ├── orderBookSnapshot.go
│   │       ├── orderEvent.go
│   │       ├── orderHistory.go
│   │       ├── orderSnapshot.go
│   │       ├── tradingPair.go
│   │       └── user.go
│   │
│   ├── connector.go
│   │
│   ├── envConfig/
│   │   └── envConfig.go
│   │
│   ├── helpers/
│   │   ├── encryption.go
│   │   ├── jwt.go
│   │   ├── request.go
│   │   └── utils.go
│   │
│   └── redis/
│       └── connect.go
│
├── migrations/
│   └── 0001_create_exchange_table.down.sql
│
├── api/
│   ├── btpin/
│   │   ├── router.go
│   │   └── service.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   └── requestValidation.go
│   │
│   ├── ngtrex/
│   │   ├── router.go
│   │   └── service.go
│   │
│   └── user/
│       ├── router.go
│       └── service.go
│
├── domain/
│   ├── balance/
│   │   ├── dto.go
│   │   └── repository.go
│   │
│   ├── exchange/
│   │   ├── btpin/
│   │   │   ├── btpin.go
│   │   │   ├── dto.go
│   │   │   └── symbols.go
│   │   │
│   │   └── nobitex/
│   │       ├── dto.go
│   │       ├── nobitex.go
│   │       └── symbols.go
│   │
│   ├── registry/
│   │   ├── exchangeFactory.go
│   │   ├── IExchange.go
│   │   └── repository.go
│   │
│   ├── exchangeCredentials/
│   │   ├── dto.go
│   │   └── repository.go
│   │
│   └── order/
│       ├── dto.go
│       └── repository.go
│
└── order/
    ├── dto.go
    ├── order.go
    └── repository.go
```

## 🏗️ How to Register and Add Modular Exchanges

Eye-On implements a **pluggable Clean Architecture** with strict layer separation between entities, service layers, and
handlers, centered around a domain-centric business layer. This design allows for seamless integration of new cryptocurrency exchanges without modifying core business logic.

### Architecture Overview

The system follows Clean Architecture principles with clear boundaries:

- **Domain Layer** (`domain/`): Contains business entities and core business rules
- **Service Layer**: Implements business logic and use cases
- **Handler Layer** (`api/`): HTTP handlers and routing logic, in which services use entity business logics.
- **Infrastructure Layer(`internals/`)**: Database, Redis, and external service connections

Each exchange is implemented as a **pluggable module** that conforms to standardized interfaces, ensuring consistency and maintainability across different exchange integrations.

### Step-by-Step Guide to Add a New Exchange

#### 1. Implement the IExchange Interface

Every exchange service **must** implement the core `IExchange` interface located in `domain/registry/IExchange.go`:

```go
type IExchange interface {
    Name() string
    Ping(ctx context.Context) error
    GetBalance(ctx context.Context, userId uuid.UUID, sign *string) ([]models.BalanceSnapshot, error)
    GetOrderBook(ctx context.Context, symbol string, userId uuid.UUID) (*models.OrderBookSnapshot, error)
    PlaceOrder(ctx context.Context, req *order.StandardOrderRequest, userId uuid.UUID) (*models.OrderHistory, error)
    CancelOrder(ctx context.Context, orderID *string, userId uuid.UUID, hours *float64) error
}
```

#### 2. Create Exchange-Specific Domain Structure

Create a new directory under `domain/exchange/` for your exchange (e.g., `domain/exchange/newexchange/`):

```
domain/exchange/newexchange/
├── newexchange.go    # Main exchange entity implementation
├── dto.go           # Data transfer objects
└── symbols.go       # Symbol registry implementation
```

#### 3. Implement Symbol Registry

Each exchange **must** have a symbol registry that implements symbol mapping and registration:

```go
type NewExchangeSymbolRegistry struct{}

func (r *NewExchangeSymbolRegistry) RegisterExchangeSymbols(exchange *models.Exchange) *[]models.TradingPair {
    // Implementation to register trading pairs specific to this exchange
    // Map exchange-specific symbols to standardized trading pairs
    tradingPairs := []models.TradingPair{
        // Define your exchange's trading pairs here
    }
    return &tradingPairs
}
```

#### 4. Register Exchange in api.go

In the `cmd/api.go` file, instantiate the symbol registry and register the exchange using the factory pattern:

```go
// Instantiate symbol registry
newExchangeSymbolRegistry := newexchange.NewExchangeSymbolRegistry{}

// Register exchange with configuration
newExchangeEntity, err := registry.GetOrCreateExchange(ctx, registry.ExchangeConfig{
    Name:          "newexchange",
    DisplayName:   "New Exchange",
    BaseURL:       "https://api.newexchange.com",
    RateLimit:     100, // requests per minute
    Features:      []string{"spot_trading", "margin_trading"}, // optional features
    SymbolFactory: &newExchangeSymbolRegistry,
})
if err != nil {
    log.Fatalf("Failed to register New Exchange: %v", err)
}
```

#### 5. Create API Service and Router with Dependency Injection

Create the service and router structure in the `api/` subdirectory:

```
api/newexchange/
├── router.go    # HTTP route definitions
└── service.go   # Business logic service layer
```

Wire up the dependencies using dependency injection:

```go
// Create router with injected dependencies
newExchangeRouter := newexchange.Router{
    Service: &newexchange.NewExchangeService{
        Exchange: &newexchangeEntity.NewExchangeExchange{
            NewExchangeExchangeModel: newExchangeEntity.Exchange,
            ExchangeRepo:             exchangeRepo,
            ExchangeCredentialRepo:   exchangeCredRepo,
            UserRepo:                 userRepo,
            TradingPairRepo:          &tradingPairRepo,
            OrderRepo:                orderRepo,
            OrderBookRepo:            orderBookRepo,
            BalanceRepo:              balanceRepo,
            Request:                  request,
        },
    },
    Parser: &jwtParser,
}

// Register routes
newExchangeRouter.RegisterRoutes(router.Group("/api/newexchange"))
```

### Key Architectural Principles

#### Clean Architecture Benefits
- **Independence**: Business logic is independent of frameworks, UI, and external services
- **Testability**: Each layer can be tested in isolation
- **Flexibility**: Easy to swap implementations without affecting business rules
- **Maintainability**: Clear separation of concerns makes code easier to understand and modify

#### Dependency Flow
The dependency flow follows the **Dependency Inversion Principle**:
- **Handlers** depend on **Services**
- **Services** depend on **Domain Entities**
- **Domain Entities** depend on **Repository Interfaces**
- **Infrastructure** implements **Repository Interfaces**

#### Business Logic Separation
- **Domain Entities**: Contains core business rules and exchange-specific logic
- **Services**: Orchestrates business operations and coordinates between entities
- **Handlers**: Handle HTTP requests/responses and input validation
- **Repositories**: Abstract data access patterns

### Best Practices for Exchange Integration

1. **Error Handling**: Implement consistent error handling across all exchange methods
2. **Rate Limiting**: Respect exchange-specific rate limits in your implementation
3. **Data Transformation**: Use DTOs to transform exchange-specific data to standardized formats
4. **Testing**: Write comprehensive unit tests for each exchange implementation
5. **Documentation**: Document exchange-specific quirks and API limitations
6. **Monitoring**: Implement logging and metrics for exchange operations

This modular approach ensures that adding new exchanges is straightforward while maintaining code quality and architectural integrity. Each exchange operates independently while adhering to the same business contracts, enabling seamless scaling and maintenance.

## 🛠️ Tech Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Containerization**: Docker & Docker Compose
- **Authentication**: JWT
- **API**: RESTful HTTP API
- **Architecture**: Clean Architecture with Pluggable Exchanges


## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Run the application**
   ```bash
   go run cmd/app.go api
   ```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:
you can use the .env.example sample for you own .env

## 📚 API Documentation

The API provides endpoints for:

- **User Management**: Registration, authentication, and exchange credentials management
- **Exchange Integration**: Bitpin, Nobitex support
- **Order Management**: Place, cancel, track orders
- **Balance Tracking**: Real-time balance updates
- **Trading Pairs**: Symbol and pair management

## 🏛️ Supported Exchanges

- **Bitpin**: Iranian cryptocurrency exchange
- **Nobitex**: Iranian cryptocurrency exchange
- **Extensible**: Easy to add new exchanges via the factory pattern

# Eye-On API Documentation 👁️

**⚠️ IMPORTANT: All endpoints use the standardized request/response format of this application, regardless of the target exchange. This ensures consistency across all exchange integrations.**

## 🔑 Authentication & Security

### ALl APi keys, access tokens etc are encrypted and stored on database.
# Bitpin Special Requirements
**🚨 CRITICAL:** Bitpin expires secret keys every 15 minutes and requires renewal using a refresh token.

- Auto-refresh functionality was not implemented due to time constraints
- A bash script has been provided for manual token refresh
- Add your refresh token to the script and execute it to update credentials
- This process must be repeated every 15 minutes to maintain active connections
- Once you got your new access token, use PUT /user/exchangeCredentials to update your bitpin credentials. 
### JWT Authentication
Most endpoints require JWT authentication via the `Authorization` header:
```
Authorization: Bearer <your_jwt_token>
```

## 📡 API Endpoints

### User Management

#### Register User
```http
POST /user/register
```

**Request Body:**
```json
{
    "username": "string (required, 3-50 chars)",
    "email": "string (required, valid email)",
    "password": "string (required, min 6 chars)"
}
```

**Response:**
```json
{
    "token": "string",
    "user": {
        "id": "uuid",
        "username": "string",
        "email": "string",
        "created_at": "timestamp"
    }
}
```

#### Login User
```http
POST /user/login
```

**Request Body:**
```json
{
    "username": "string (required)",
    "password": "string (required)"
}
```

**Response:**
```json
{
    "token": "string",
    "user": {
        "id": "uuid",
        "username": "string",
        "email": "string"
    }
}
```

#### Create Exchange Credentials
## All Credential Tokens are encrypted before storing on the database
```http
POST /user/exchangeCredentials
Authorization: Bearer <token>
```

**Request Body:**
```json
{
    "exchange_name": "string (required)",
    "label": "string (required)",
    "api_key": "string (required)",
    "secret_key": "string (required)",
    "access_key": "string (optional)",
    "is_testnet": "boolean"
}
```

#### Update Exchange Credentials
```http
PUT /user/exchangeCredentials
Authorization: Bearer <token>
```

**Request Body:**
```json
{
    "exchange_name": "string (required)",
    "label": "string (required)",
    "api_key": "string (required)",
    "secret_key": "string (required)",
    "access_key": "string (optional)",
    "is_testnet": "boolean"
}
```

### Exchange Operations

All exchange endpoints follow the same standardized format regardless of the target exchange (Bitpin, Nobitex, etc.).

#### Place Order
## Be Aware that Nobitex exchange only supports transactions where base currency is either tether or rials.
```http
POST /exchange/{exchange_name}/order
Authorization: Bearer <token>
```

**Standard Request Body:**
### it can be more consistent and could have been only one request without caring about which one is required
### i just didin`t find enough time to do that
```json
{
    "symbol": "string (required)",
    "side": "buy|sell (required)",
    "type": "market|limit (required)",
    "quantity": "number (optional)",
    "base_currency": "string (only required for nobitex)",
    "quote_currency": "string (only required for nobitex)",
    "base_amount": "number (only required for bitpin)",
    "quote_amount": "number (only required for bitpin)",
    "price": "number (only required for nobitex)",
    "stop_price": "number (optional)",
    "time_in_force": "string (optional)",
    "client_order_id": "string (optional)"
}
```

**Exchange-Specific Requirements:**

**Bitpin Requirements:**
- `symbol` (required)
- `side` (required)
- `type` (required)
- `base_amount` (required)
- `quote_amount` (required)
- `price` (required)
- `client_order_id` (required)

**Nobitex Requirements:**
- `symbol` (required)
- `type` (required)
- `price` (required for limit orders)
- `base_currency` (required)
- `quote_currency` (required)
- `quantity` (required)

#### Cancel Order
```http
DELETE /exchange/{exchange_name}/order/{orderId}
Authorization: Bearer <token>
```

**Request Body:**
```json
{
    "orderId": "string (required)",
    "hours": "number (optional, required for Nobitex)"
}
```

**Note:** The `hours` parameter is only required for Nobitex exchange operations.

#### Get Order Book
```http
GET /exchange/{exchange_name}/orderBook/{symbol}
Authorization: Bearer <token>
```

**Path Parameters:**
- `symbol`: Trading pair symbol (e.g., "BTCUSDT")

#### Get Balance
```http
GET /exchange/{exchange_name}/balance
Authorization: Bearer <token>
```

## 🔧 Exchange-Specific Endpoints

### Bitpin
```http
POST /exchange/bitpin/order
DELETE /exchange/bitpin/order/{orderId}
GET /exchange/bitpin/orderBook/{symbol}
GET /exchange/bitpin/balance
```

### Nobitex
```http
POST /exchange/nobitex/order
DELETE /exchange/nobitex/order/{orderId}
GET /exchange/nobitex/orderBook/{symbol}
GET /exchange/nobitex/balance
```

## 📊 Standard Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
    "success": true,
    "data": {
        // Response data specific to the endpoint
    },
    "message": "Operation completed successfully"
}
```

**Error Response:**
```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "Human readable error message",
        "details": "Additional error details"
    }
}
```

## 🔐 Security Considerations

1. **Token Management**: JWT tokens should be stored securely and refreshed as needed
2. **Bitpin Token Refresh**: Implement automated refresh mechanism for production use
3. **Rate Limiting**: Respect exchange-specific rate limits
4. **Credential Security**: Exchange credentials are encrypted at rest
5. **HTTPS Only**: All API calls should use HTTPS in production

## 🚀 Getting Started

1. Register a user account
2. Login to receive JWT token
3. Add exchange credentials for your preferred exchange
4. For Bitpin: Set up token refresh mechanism
5. Start trading using the standardized API endpoints

## 📋 Order Types & Sides

**Order Sides:**
- `buy`: Purchase order
- `sell`: Sale order

**Order Types:**
- `market`: Execute immediately at current market price
- `limit`: Execute only at specified price or better

**Time in Force Options:**
- `GTC`: Good Till Canceled
- `IOC`: Immediate Or Cancel
- `FOK`: Fill Or Kill

## ⚠️ Important Notes

- All monetary values should be provided as numbers, not strings
- Symbol formats may vary between exchanges - check exchange documentation
- Order IDs are exchange-specific and should be stored for tracking
- Some exchanges may have minimum order quantities or amounts
- Always validate order parameters before submission to avoid failed transactions
