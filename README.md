# Eye-On 👁️

**Exchange and vendor agnostic transaction manager**

A sophisticated Go-based trading and exchange management system that provides unified access to multiple cryptocurrency exchanges with robust order management, user authentication, and real-time data processing.

---

## 🚀 Features

- **Multi-Exchange Support**: Unified interface for multiple cryptocurrency exchanges (Bitpin, Nobitex)
- **Order Management**: Comprehensive order tracking, history, and snapshot capabilities
- **User Authentication**: JWT-based authentication with secure credential management
- **RESTful API**: Clean HTTP API with middleware for authentication and validation
- **Database Migrations**: Structured database schema management
- **Docker Ready**: Complete containerization with Docker Compose
- **Bitpin Auto-Refresh**: Automated refresh of Bitpin access tokens using refresh tokens
- **Standardized Symbols**: Nobitex exchange symbol bug fixed; symbols are now standardized

---

## 🏗️ Architecture

The project follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Business logic and entities
- **API Layer**: HTTP handlers and routing
- **Infrastructure Layer**: Database, Redis, and external service connections
- **Application Layer**: Use cases and application services

---

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
│   ├── envConfig/
│   ├── helpers/
│   └── redis/
│
├── migrations/
│
├── api/
│   ├── btpin/
│   ├── middleware/
│   ├── ngtrex/
│   └── user/
│
├── domain/
│   ├── balance/
│   ├── exchange/
│   ├── registry/
│   ├── exchangeCredentials/
│   └── order/
│
└── order/
```

---

## 🧩 Recent Updates

* ✅ **Bitpin Token Auto-Renewal Implemented**  
  Bitpin access tokens are now **automatically refreshed** using the refresh token. Manual renewal using a script and updating via `PUT /user/exchangeCredentials` is no longer necessary.

* 🛠️ **Nobitex Symbol Standardization Fix**  
  A bug in the Nobitex exchange symbol mapping has been fixed. The symbol registry now uses a **consistent and correct standard format**.

* 🐳 **Docker Deployment Ready**  
  A `Dockerfile` and `docker-compose.yml` have been added for **seamless containerized deployment**.

---

## 🛠️ Tech Stack

* **Backend**: Go 1.21+
* **Database**: PostgreSQL 15
* **Cache**: Redis 7
* **Containerization**: Docker & Docker Compose
* **Authentication**: JWT
* **API**: RESTful HTTP API
* **Architecture**: Clean Architecture with Pluggable Exchanges

---

## 🚀 Quick Start

### Prerequisites

* Docker and Docker Compose
* Go 1.21+ (for local development)

### Local Development

```bash
go mod download
go run cmd/app.go api
```

### With Docker Compose

```bash
docker-compose -f infra/docker-compose.yml up --build
```

---

## 🔧 Configuration

Create a `.env` file in the root directory. Use `.env.example` as a starting point.

---

## 📚 API Documentation

Provides endpoints for:

* **User Management**
* **Exchange Integration** (Bitpin, Nobitex)
* **Order Placement & Tracking**
* **Balance Retrieval**
* **Trading Symbol Mapping**

---

## 🏛️ Supported Exchanges

* **Bitpin** – Access token refresh fully automated
* **Nobitex** – Symbol bug fixed and standardized
* **Extensible** – New exchanges pluggable via a factory interface

---

## 🔑 Authentication

All API calls (except registration/login) require JWT tokens in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

---

## 🔐 Security Notes

* All exchange credentials are encrypted before being stored
* JWT tokens are securely generated and must be protected
* Bitpin tokens are auto-refreshed in the background
* Always use HTTPS in production deployments

---

## 📡 Example API Endpoints

### User Registration

```http
POST /user/register
```

### User Login

```http
POST /user/login
```

### Create Exchange Credentials

```http
POST /user/exchangeCredentials
```

### Update Exchange Credentials

```http
PUT /user/exchangeCredentials
```

### Place Order

```http
POST /exchange/{exchange_name}/order
```

### Cancel Order

```http
DELETE /exchange/{exchange_name}/order/{orderId}
```

### Get Balance

```http
GET /exchange/{exchange_name}/balance
```

### Get Order Book

```http
GET /exchange/{exchange_name}/orderBook/{symbol}
```

---

## 📋 Order Format Overview

**Standard Fields:**

* `symbol`: string
* `side`: buy/sell
* `type`: market/limit
* `price`, `quantity`, `base_currency`, `quote_currency`, `base_amount`, etc., depending on exchange

**Exchange-Specific Notes:**

* Nobitex supports only Tether and IRR as base currencies
* Bitpin requires both `base_amount` and `quote_amount`

---

## 🧠 Design Philosophy

* Clean Architecture principles
* Strong separation of concerns
* Each exchange implementation is modular and independent
* Easily testable and maintainable
* Secure handling of all user and exchange credentials

---

## 🚧 Contribution Notes

When adding a new exchange:

1. Implement the `IExchange` interface
2. Create exchange domain structure
3. Register via factory in `cmd/api.go`
4. Implement router and service logic in `api/{exchange}`
5. Add symbol mapping logic and ensure standardization

---

## ✅ Summary

Eye-On offers a clean, pluggable, and scalable way to integrate multiple crypto exchanges. With the **automated Bitpin token refresh**, **Nobitex symbol standardization fix**, and **Dockerized deployment**, it is production-ready for modern crypto trading APIs.