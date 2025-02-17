package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/schollz/progressbar/v3"
)

type Record struct {
	ID      string    `json:"id"`
	Date    time.Time `json:"date"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

func main() {
	table := flag.String("table", "", "Destination table (redis_engine_table or kafka_engine_table)")
	interval := flag.Duration("interval", 100*time.Millisecond, "Interval between insertions")
	count := flag.Int("count", 1_000, "Number of records to insert")
	flag.Parse()

	if *table == "" {
		log.Fatal("Please specify a destination table using -table flag")
	}

	validTables := map[string]bool{
		"redis_null": true,
		"kafka_null": true,
	}

	if !validTables[*table] {
		log.Fatalf("Invalid table name. Valid options are: %s", strings.Join(getKeys(validTables), ", "))
	}

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

	if conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	faker := gofakeit.New(0)

	logFileName := fmt.Sprintf("insert_%s.log", *table)
	lf, err := os.Create(logFileName)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer lf.Close()

	// Create a progress bar
	bar := progressbar.Default(int64(*count))

	// Generate and insert data in a loop
	for i := 0; i < *count; i++ {
		record := generateRecord(faker)

		// Insert into specified table
		err = insertRecord(conn, *table, record, lf)
		if err != nil {
			log.Fatalf("Error inserting into %s: %v", *table, err)
		}

		// Print the generated record
		// jsonData, _ := json.MarshalIndent(record, "", "  ")
		// fmt.Printf("Generated record for %s:\n%s\n", *table, string(jsonData))

		bar.Add(1)
		time.Sleep(*interval)
	}
	fmt.Printf("Done, check %s for logs\n", logFileName)
}
func generateRecord(faker *gofakeit.Faker) Record {

	return Record{
		ID:      generateTimestampID(),
		Date:    time.Now().UTC(),
		Email:   faker.Email(),
		Message: faker.HackerPhrase(),
	}
}

func generateTimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func insertRecord(conn clickhouse.Conn, table string, record Record, lf *os.File) error {
	start := time.Now()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			ID,
			Date,
			Email,
			Message
		) VALUES (?, ?, ?, ?)
	`, table)

	err := conn.Exec(context.Background(), query,
		record.ID,
		record.Date,
		record.Email,
		record.Message,
	)
	end := time.Now()
	duration := end.Sub(start)
	fmt.Fprintf(lf, "%s,%v\n", table, duration.Milliseconds())
	return err
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
