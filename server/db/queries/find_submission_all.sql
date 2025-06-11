SELECT id, TIMESTAMP, task_id, score, COUNT(*) OVER() AS total_amount
FROM Submission
WHERE owner_name = ?
ORDER BY TIMESTAMP DESC
LIMIT ?;