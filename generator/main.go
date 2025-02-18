package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
)

type Record struct {
	ID      string    `json:"id"`
	Date    time.Time `json:"date"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

type Config struct {
	Table      string
	Interval   time.Duration
	BatchCount int
	BatchSize  int
}

func loadConfig() Config {
	config := Config{}

	flag.StringVar(&config.Table, "table", "", "Destination table (redis_engine_table, kafka_engine_table, postgres_engine_table)")
	flag.DurationVar(&config.Interval, "interval", 70*time.Millisecond, "Interval between insertions")
	flag.IntVar(&config.BatchCount, "batch-count", 10, "Number of batches to insert")
	flag.IntVar(&config.BatchSize, "batch-size", 7, "Number of records per batch")
	flag.Parse()

	if table := os.Getenv("CLICKHOUSE_TABLE"); table != "" {
		config.Table = table
	}
	if interval := os.Getenv("CLICKHOUSE_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			config.Interval = d
		}
	}
	if batchCount := os.Getenv("CLICKHOUSE_BATCH_COUNT"); batchCount != "" {
		if n, err := strconv.Atoi(batchCount); err == nil {
			config.BatchCount = n
		}
	}
	if batchSize := os.Getenv("CLICKHOUSE_BATCH_SIZE"); batchSize != "" {
		if n, err := strconv.Atoi(batchSize); err == nil {
			config.BatchSize = n
		}
	}

	return config
}

func main() {
	config := loadConfig()

	if config.Table == "" {
		log.Fatal("Please specify a destination table using -table flag or CLICKHOUSE_TABLE env var")
	}

	validTables := map[string]bool{
		"redis_null":    true,
		"kafka_null":    true,
		"postgres_null": true,
	}

	if !validTables[config.Table] {
		log.Fatalf("Invalid table name. Valid options are: %s", strings.Join(getKeys(validTables), ", "))
	}

	totalRecords := config.BatchCount * config.BatchSize
	log.Printf("Will insert %d records in %d batches of %d records each",
		totalRecords, config.BatchCount, config.BatchSize)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	faker := gofakeit.New(0)

	logFileName := fmt.Sprintf("insert_%s.log", config.Table)
	lf, err := os.Create(logFileName)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer lf.Close()

	bar := progressbar.Default(int64(totalRecords))

	for batchNum := 0; batchNum < config.BatchCount; batchNum++ {
		start := time.Now()

		// Prepare batch
		batch, err := conn.PrepareBatch(ctx, fmt.Sprintf("INSERT INTO %s", config.Table))
		if err != nil {
			log.Fatalf("Error preparing batch: %v", err)
		}

		// Generate and append records
		for i := 0; i < config.BatchSize; i++ {
			record := generateRecord(faker)
			err := batch.Append(
				record.ID,
				record.Date,
				record.Email,
				record.Message,
			)
			if err != nil {
				log.Fatalf("Error appending to batch: %v", err)
			}
			bar.Add(1)
		}

		// Send batch
		err = batch.Send()
		if err != nil {
			log.Fatalf("Error sending batch: %v", err)
		}

		duration := time.Since(start)
		fmt.Fprintf(lf, "%s,%v,%d\n", config.Table, duration.Milliseconds(), config.BatchSize)

		time.Sleep(config.Interval)
	}

	fmt.Printf("\nDone, check %s for logs\n", logFileName)
}

func generateRecord(faker *gofakeit.Faker) Record {
	return Record{
		ID:      uuid.New().String(),
		Date:    time.Now().UTC(),
		Email:   faker.Email(),
		Message: faker.HackerPhrase(),
	}
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
