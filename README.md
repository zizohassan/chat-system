# Chat System with Go, Cassandra, Redis, Prometheus, Grafana, and Loki

## Overview

This project is a simplified chat platform built using Go. It handles user authentication, message sending, and message
retrieval. The backend uses Cassandra for data storage, Redis for caching, Prometheus for metrics, Grafana for
visualization, and Loki for log aggregation.

## Features

- User Registration and Login
- Send and Retrieve Messages
- Distributed Data Storage with Cassandra
- Caching with Redis
- Monitoring with Prometheus
- Visualization with Grafana
- Log Aggregation with Loki

## Prerequisites

- Docker
- Docker Compose

## Setup

### Step 1: Create Environment File

Create a `.env` file in the root directory with the following command:

``` 
cp .env.example .env
```

### Step 2: Start Docker Compose

Build and start all services using Docker Compose:

```
docker-compose down --rmi all
docker-compose up --build
```

# Usage

## Accessing Services

- Grafana: http://localhost:3000
    - Default credentials: admin / admin
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
- Nginx (API Gateway): http://localhost:5050
    - Access Swagger documentation at http://localhost:5050/swagger

## Monitoring and Logs

### Prometheus
Prometheus is set up to collect metrics from your Go application. <br>
Access it at http://localhost:9090.

### Grafana
Grafana is configured with dashboards to visualize metrics and logs. <br>
Access it at http://localhost:3000. 
Default credentials: 
``` admin / admin ```

### Loki
Loki is used for log aggregation and can be queried within Grafana. <br>
Ensure Loki data source is configured in Grafana.

### Additional Information

#### Docker Compose Configuration
The ```docker-compose.yml``` file sets up the entire stack including Cassandra, Redis, Prometheus, Grafana, Loki, and the Go
application. It also includes health checks and volume mounts to ensure data persistence.

#### Loki Configuration
Loki's configuration is specified in ```/docker/loki/local-config.yaml``` and mounted in the Docker Compose setup.

#### Grafana Provisioning
Grafana is provisioned with data sources and dashboards via configuration files in ```grafana/provisioning/```


## Architectural Decisions and Assumptions

### Architectural Decisions

- **Microservices Architecture:** The chat application is built as a set of microservices to ensure scalability and maintainability.
- **Cassandra for Data Storage:** Cassandra is chosen for its distributed nature, providing high availability and scalability.
- **Redis for Caching:** Redis is used as an in-memory cache to improve the performance of frequently accessed data.
- **Docker and Docker Compose:** Docker is used to containerize the application, and Docker Compose is used to manage multi-container deployments.
- **Prometheus for Monitoring:** Prometheus is used for monitoring system metrics and alerting.
- **Grafana for Visualization:** Grafana is used for visualizing metrics and logs collected by Prometheus and Loki.
- **Loki for Log Aggregation:** Loki is used to aggregate logs from the application, making it easier to search and analyze logs.

### Reflex Documentation
Reflex is used for watching file changes and triggering rebuilds during development. This helps streamline the development workflow by automatically recompiling the Go application when changes are detected.
To learn more about Reflex, refer to the Reflex Documentation.

### Assumptions

- **Single Node Deployment:** The setup assumes a single node deployment for simplicity. In a production environment, Cassandra and other services should be deployed in a distributed manner.
- **Basic Authentication:** The application uses basic username-password authentication. In a real-world application, OAuth or other advanced authentication mechanisms might be used.
- **Default Credentials for Grafana:** The default credentials for Grafana are set to admin / admin. It's recommended to change these credentials in a production setup.
- **Development Environment:** The setup is intended for a development environment. Additional security and performance optimizations would be needed for production.
