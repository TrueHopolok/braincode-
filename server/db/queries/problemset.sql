SELECT id, title, COUNT(*) OVER(PARTITION BY id) AS task_amount
FROM task
LIMIT ? OFFSET ?;