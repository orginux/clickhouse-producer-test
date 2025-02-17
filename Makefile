up:
	@docker compose up -d

wait-for-clickhouse:
	@docker exec clickhouse-producer-test-clickhouse-1 bash -c "until wget --no-verbose --tries=1 --spider localhost:8123/ping; do sleep 1; done"

down:
	@docker compose down

clean:
	docker compose down -v

copy-schema:
	@docker cp schema.sql clickhouse-producer-test-clickhouse-1:/tmp/schema.sql

apply-schema: copy-schema
	@docker exec -i clickhouse-producer-test-clickhouse-1 clickhouse-client --queries-file /tmp/schema.sql

go-tidy:
	@cd generator && go mod tidy

build-generator: go-tidy
	@cd generator && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/generator
	@du -sh generator/bin/generator

copy-generator-to-container: build-generator
	@docker cp generator/bin/generator clickhouse-producer-test-clickhouse-1:/usr/local/bin

start-generator-redis: copy-generator-to-container
	@docker exec clickhouse-producer-test-clickhouse-1 /usr/local/bin/generator -table redis_null

start-generator-kafka: copy-generator-to-container
	@docker exec clickhouse-producer-test-clickhouse-1 /usr/local/bin/generator -table kafka_null

copy-query-to-container:
	@docker cp query.sql clickhouse-producer-test-clickhouse-1:/tmp/query.sql

run-query: copy-query-to-container
	@docker exec clickhouse-producer-test-clickhouse-1 clickhouse-local --queries-file /tmp/query.sql

test: clean up wait-for-clickhouse apply-schema start-generator-redis start-generator-kafka run-query down

clickhouse-client:
	@docker exec -it clickhouse-producer-test-clickhouse-1 clickhouse-client -n

