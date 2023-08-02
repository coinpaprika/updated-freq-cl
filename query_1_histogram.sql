-- This query will show the distribution of the time between price changes as a histogram
WITH diff_between_price_changes_based_on_now_timestamp
         AS (SELECT now_timestamp - lagInFrame(now_timestamp) OVER (ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) AS diff
             FROM file('result.csv')
             WHERE price_changed = 1 ORDER BY now_timestamp ASC)
SELECT bucket * 10 diff_secs, count(*) count, bar(count, 0, 43, 10) hitogram
FROM diff_between_price_changes_based_on_now_timestamp GROUP BY intDiv(diff, 10) as bucket;