# Go Application Task

## Overview

This Go application provides a simple API for order management. The application allows users to:

- **Login using JWT** for authentication.
- **Create Orders**.
- **Show Orders**.
- **Cancel Orders**.

## Prerequisites

- Docker
- Go (1.21 or higher)
- PostgreSQL (Database for storing orders)

## Setup Instructions

### 1. Clone the Repository
### 2. Setup Environment
### 3. Build with `go build -o cmd/main .`
### 4. Run with `go build -o cmd/main .`

By default, the API will be available at http://localhost:8180

Alternatively, if you're using Docker, you can run the app using the provided docker-compose.yml file:

```bash
docker-compose up