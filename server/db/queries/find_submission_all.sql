SELECT id, TIMESTAMP, task_id, score, COUNT(*) OVER() AS total_amount
FROM submission
WHERE owner_name = ?
ORDER BY TIMESTAMP;