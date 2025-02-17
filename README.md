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

Start all services and run the complete test:
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

The data generator supports the following configuration options:

```bash
# Command line flags
-table string        # Target table (redis_null or kafka_null)
-interval duration   # Delay between batch insertions (default: 200ms)
-batch-count int    # Number of batches to insert (default: 100)
-batch-size int     # Number of records per batch (default: 10)

# Environment variables
CLICKHOUSE_TABLE        # Same as -table
CLICKHOUSE_INTERVAL     # Same as -interval
CLICKHOUSE_BATCH_COUNT  # Same as -batch-count
CLICKHOUSE_BATCH_SIZE   # Same as -batch-size
```

### Batch Processing
The generator uses ClickHouse's native batch processing capabilities. The total number of records inserted will be:
```
total_records = batch_count × batch_size
```

For example:
```bash
# Insert 1000 records in 100 batches of 10 records each
./generator -table redis_null -batch-count 100 -batch-size 10

# Insert 10000 records in 100 batches of 100 records each
./generator -table kafka_null -batch-count 100 -batch-size 100
```

Each batch insertion's performance is logged separately, allowing for detailed analysis of batch processing efficiency.

## Dependencies
### Go packages:
- github.com/ClickHouse/clickhouse-go/v2
- github.com/brianvoe/gofakeit/v7
- github.com/schollz/progressbar/v3
