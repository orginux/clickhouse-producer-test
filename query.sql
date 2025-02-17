SELECT
        c1 as table_name,
        max(c2) as max_execution_ms,
        min(c2) as min_execution_ms,
        round(avg(c2),2) as avg_execution_ms,
        round(quantile(0.5)(c2),2) as median_execution_ms,
        round(quantile(0.75)(c2),2) as p75_execution_ms,
        round(quantile(0.95)(c2), 2) as p95_execution_ms,
        round(quantile(0.99)(c2), 2) as p99_execution_ms,
        sum(c3) as total_rows,
        count() as total_inserts
FROM 'insert_*.log'
GROUP BY table_name
FORMAT PrettyNoEscapesMonoBlock
;
