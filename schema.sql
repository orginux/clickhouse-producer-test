DROP TABLE IF EXISTS kafka_engine_table;
CREATE TABLE kafka_engine_table
(
        ID String,
        Date DateTime,
        Email String,
        Message String
)
ENGINE = Kafka
SETTINGS
        kafka_broker_list = 'kafka:29092',
        kafka_topic_list = 'topic1',
        kafka_group_name = 'clickhouse_group',
        kafka_format = 'JSONEachRow'
;

DROP TABLE IF EXISTS kafka_null;
CREATE TABLE kafka_null AS kafka_engine_table
ENGINE = Null
;

CREATE MATERIALIZED VIEW kafka_mv TO kafka_engine_table
AS SELECT * FROM kafka_null;

-- 
DROP TABLE IF EXISTS redis_engine_table;
CREATE TABLE redis_engine_table
(
        ID String,
        Date DateTime,
        Email String,
        Message String
)
ENGINE = Redis('redis:6379')
PRIMARY KEY (ID)
;

DROP TABLE IF EXISTS redis_null;
CREATE TABLE redis_null AS redis_engine_table
ENGINE = Null
;

CREATE MATERIALIZED VIEW redis_mv TO redis_engine_table
AS SELECT
        toString(ID),
        Date,
        Email,
        Message
FROM redis_null;
