SELECT id, title, COUNT(*) OVER() AS totalAmount
FROM task
LIMIT ? OFFSET ?;