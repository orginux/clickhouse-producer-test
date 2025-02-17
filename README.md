# ClickHouse Producer Test

A project to test ClickHouse's capabilities as a data producer for Redis and Kafka, measuring insert performance and data throughput.

## Project Goal

The main purpose of this project is to evaluate ClickHouse's performance as a data producer, specifically.

## Project Structure

```
.
├── docker-compose.yaml # Infrastructure setup
├── Makefile            # Command automation
├── generator/          # Data generator source code in Go
├── schema.sql          # ClickHouse tables definitions
└── query.sql           # Performance queries
```

## Quick Start

1. Start all services and run the complete test:
```bash
make test
```

This will:
- Clean up any existing containers
- Start all services
- Wait for ClickHouse to be ready
- Apply the database schema
- Run generators for both Redis and Kafka
- Generate performance report
- Clean up

## Available Make Commands

```bash
# Infrastructure
make up                    # Start all containers
make down                  # Stop all containers
make clean                # Remove containers and volumes

# Schema Management
make apply-schema         # Apply ClickHouse schema

# Generator
make build-generator      # Build the data generator
make start-generator-redis # Run generator for Redis
make start-generator-kafka # Run generator for Kafka

# Monitoring
make run-query           # Run performance analysis

# Development
make clickhouse-client   # Open ClickHouse CLI
make go-tidy            # Update Go dependencies
```

## ClickHouse Schema

The test data includes:
```sql
CREATE TABLE kafka_engine_table (
    ID String,
    Date DateTime,
    Email String,
    Message String
) ENGINE = Kafka ...

CREATE TABLE redis_engine_table (
    ID String,
    Date DateTime,
    Email String,
    Message String
) ENGINE = Redis ...
```

## Performance Analysis

Performance metrics are logged to files named `insert_{table}.log`. Analysis query provides:
- Maximum execution time
- Minimum execution time
- Median execution time
- 75th percentile
- 95th percentile
- 99th percentile
- Total number of queries

To view performance metrics:
```bash
make run-query
```

## Generator Configuration

The data generator supports:
```bash
-table string        # Target table (redis_null or kafka_null)
-interval duration   # Delay between insertions (default: 100ms)
-count int          # Number of records to insert (default: 1000)
```

## Dependencies
### Go packages:
- github.com/ClickHouse/clickhouse-go/v2
- github.com/brianvoe/gofakeit/v7
- github.com/schollz/progressbar/v3
