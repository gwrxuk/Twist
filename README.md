# Twist - Blockchain Infrastructure Management Platform

A high-performance, scalable platform for managing and optimizing EVM-compatible blockchain nodes and infrastructure.

## Overview

Twist is a comprehensive platform designed for blockchain engineers and infrastructure teams to manage, monitor, and optimize blockchain nodes across multiple EVM-compatible chains. The platform enables seamless deployment, configuration, and health monitoring of blockchain nodes in both cloud and on-premise environments.

## Key Features

- **Node Management**: Setup, configure, and monitor EVM blockchain nodes
- **Performance Optimization**: Tools to optimize blockchain node performance
- **Smart Contract Analysis**: Analyze and optimize Solidity smart contracts
- **Cloud Integration**: Seamless deployment to AWS and Google Cloud
- **Multi-chain Support**: Support for Ethereum, Polygon, Arbitrum, BSC, and other EVM chains

## Project Structure

The project is organized as a microservices architecture with the following components:

```
twist/
├── core-engine/             # Rust-based high-performance node management engine
├── api-gateway/             # Go-based API gateway for external access
├── smart-contract-service/  # Solidity smart contracts and deployment tools
├── dashboard/               # TypeScript/Next.js frontend (to be implemented)
├── analytics-service/       # Rust-based analytics engine (to be implemented)
├── monitoring/              # Prometheus and Grafana configuration
└── docker-compose.yml       # Docker Compose configuration for local development
```

## Architecture

Twist is built using a microservices architecture:

- **Core Engine** (Rust): High-performance node management and optimization
- **API Gateway** (Go): RESTful and gRPC API interfaces
- **Dashboard** (TypeScript/Next.js): User interface for monitoring and management
- **Smart Contract Service** (Go/Solidity): Smart contract deployment and analysis
- **Analytics Service** (Rust): Performance analytics and optimization recommendations

## Technology Stack

- **Languages**: Rust, Go, TypeScript, Solidity
- **Frameworks**: TRPC, ExpressJS, Next.js, Actix-web
- **Infrastructure**: Docker, Kubernetes, AWS, Google Cloud
- **Blockchain**: Ethereum, Polygon, BSC, Arbitrum
- **Databases**: PostgreSQL, Redis
- **Monitoring**: Prometheus, Grafana

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Rust 1.70 or later
- Go 1.20 or later
- Node.js 18 or later
- PostgreSQL 15
- Redis 7

### Environment Setup

1. Clone the repository
```bash
git clone https://github.com/your-org/twist.git
cd twist
```

2. Create environment files
```bash
cp core-engine/example.env core-engine/.env
cp api-gateway/example.env api-gateway/.env
cp smart-contract-service/env.example smart-contract-service/.env
```

3. Update the environment files with your configuration

### Running with Docker

Start all services using Docker Compose:

```bash
docker-compose up -d
```

### Running Services Individually

#### Core Engine (Rust)

```bash
cd core-engine
cargo run
```

#### API Gateway (Go)

```bash
cd api-gateway
go run cmd/main.go
```

#### Smart Contract Service

```bash
cd smart-contract-service
npm install
npx hardhat compile
npx hardhat test
```

### Accessing Services

- Core Engine API: http://localhost:8080
- API Gateway: http://localhost:8000
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

## Development

### Core Engine

The Core Engine is built with Rust and provides high-performance blockchain node management. To build:

```bash
cd core-engine
cargo build --release
```

### API Gateway

The API Gateway is built with Go and provides external access to the platform. To build:

```bash
cd api-gateway
go build -o api-gateway cmd/main.go
```

### Smart Contracts

The platform includes Solidity smart contracts for on-chain node registry and token governance:

```bash
cd smart-contract-service
npx hardhat compile
npx hardhat test
```

## Deployment

### Cloud Deployment

The platform can be deployed to AWS or Google Cloud using Kubernetes. Deployment scripts and Helm charts will be provided in future updates.

### On-Premise Deployment

For on-premise deployment, you can use Docker Compose or Kubernetes. Ensure that your environment meets the system requirements.

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

## License

MIT

## Contact

For questions or support, please contact [team@twist.example.com](mailto:team@twist.example.com). 