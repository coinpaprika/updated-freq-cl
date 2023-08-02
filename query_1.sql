WITH diff_between_price_changes_based_on_now_timestamp
         AS (SELECT now_timestamp - lagInFrame(now_timestamp) OVER (ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) AS diff
             FROM file('result.csv')
             WHERE price_changed = 1 ORDER BY now_timestamp ASC)
SELECT max(diff) as max_diff, avg(diff) as avg_diff, stddevPopStable(diff) as stddev_diff
FROM diff_between_price_changes_based_on_now_timestamp;