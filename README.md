# go ddd example

## Overview

Thank you for the opportunity, as per our discussion, i have implemented the service by using golang and DDD.
This service is handling payment requests for multiple sold items from multiple sellers by respecting business transaction limit.
Service includes REST API interface.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional)

## Setup

### Install Dependencies

```
go mod tidy
```

- with make:

```
make tidy
```

### Configure Environment Variables

```
PG_HOST=localhost
PG_USER=vcpayout
PG_PASSWORD=vcpayout
PG_DB_NAME="vcpayout"
PG_PORT=5432

```

You can pass this step, since it is a test project i am sending the .env file as well

## Runing the project

```
go run ./cmd/restapi/main.go
```

- with make:
  make will handle postgres docker instance and migration and seeding

```
make run
```

### Run Database Migration (if make is not used)

migration/main.go will migrate the tables also will create some item seeds.

```

docker-compose -f ./docker/docker-compose.yaml --env-file ./.env up postgres -d

go run ./cmd/migrate/main.go

```

- with make:
  this will automatically provision the docker instance and migrate

```
make migrate
```

## API Documentation

API docs available once you run the project at:

- <http://localhost:3000/swagger/index.html>

## Running Tests

Due to time constraints i implemented as much as possible unit tests. However did not add the integration tests.
For integration test i would use testcontainers package.

- <https://testcontainers.com/?language=go>

### Running unit tests

```
go test ./...
```

- with make

```
make test
```

## Solution

When items are requested to be paid to client (assuming sold items) the flow works like below;

1. We get the transaction limit (this one is mocked)
2. We are converting the transaction limit to requested currency.As example lets say, if the payout requested USD, but transaction limit is 100 GBP, transaction limit should reflect the USD.
3. Since we want to transfer as much money as possible within the transaction limit (to my understanding), instead of calculating and paying out each item we aggregate them and calculating total amount of requested items. As example if someone sold some items valued 200 GBP and 400 GBP and if transaction limit is 100 GBP, we batch them and create 5 100GBP payouts. Not implemented but later on it would be possible to pay first one since it is within the transaction limit and schedule and spread other ones.
4. We create a sale record to the item referancing itemID and batchPayoutId

## Project Structure

```
.
├── cmd
│   ├── migrate
│   └── restapi
├── config
├── docker
├── docs
├── internal
│   ├── application
│   ├── common
│   ├── domain
│   ├── infrastructure
│   └── interface
├── migrations
└── README.md

```

### Interface Layer

Interface layer is responsible for handling the user requests. At the moment we only rest api however in feature
this could be extended with grpc, graphql, websocket, message queues etc.

### Infastructure Layer

Infastructure layer is responsible for external dependencies like databases, caching, messaging queues, 3rd party interfaces etc. It acts as a bridge between application layer and
external services.

### Application Layer

Application layer is responsible for orchestrating the flow of data between user facing interfaces, this layer does not contain business logic, the whole responsibility is cordinating tasks
from infrastructure layer and domain layer.

### Domain Layer

Domain layer contains core business logic and business rules. This layer is independent of the external systems.

# Conclusion

Thank you for reviewing this project, as per our discussion with Kamal, i did implemented using golang however if it is requested i can implement with typescript as well.
There are many areas to improve such as:

- Implementing integration tests
- Use of domain events and event sourcing
- Observability
- Cleaner logging
