BEGIN;

-- Dates: last 365 days from CURRENT_DATE (rolling window). Guided Investigate
-- scenarios read live MIN(date)/MAX(date) and inject DATE literals so sample
-- SQL returns rows without freezing a calendar year.

CREATE SCHEMA IF NOT EXISTS demo;

INSERT INTO demo.sales (
    id,
    date,
    product_category,
    product_name,
    quantity,
    unit_price,
    total_amount,
    region,
    sales_rep
)
SELECT
    gen_random_uuid(),
    (CURRENT_DATE - (random() * 365)::int),
    (ARRAY['Electronics','Furniture','Office Supplies','Clothing','Accessories'])[1 + (random() * 4)::int],
    (ARRAY['Alpha','Beta','Gamma','Delta','Epsilon','Zeta'])[1 + (random() * 5)::int],
    1 + (random() * 20)::int,
    ROUND((10 + random() * 490)::numeric, 2),
    ROUND((10 + random() * 490)::numeric, 2) * (1 + (random() * 20)::int),
    (ARRAY['North','South','East','West','Central'])[1 + (random() * 4)::int],
    (ARRAY['A. Lee','B. Singh','C. Patel','D. Kim','E. Garcia'])[1 + (random() * 4)::int]
-- 300k rows across the monthly partitions (~55 MB, ~3s to insert).
--
-- 8000 rows was too small to demonstrate anything: at ~160 rows per partition
-- the whole table sits in shared buffers, partition pruning saves microseconds,
-- and run-to-run timing noise exceeded the difference being shown — the same
-- query measured 2ms and 6ms on consecutive runs. At this size the demo
-- comparison reports 24ms -> 7ms, which is a difference a reader can trust.
--
-- For the 10M-row benchmark used in the case study, use `make seed-large-docker`.
FROM generate_series(1, 300000);

COMMIT;
