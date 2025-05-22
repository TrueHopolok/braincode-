SELECT id, title, owner_name, COUNT(*) OVER() AS totalAmount
FROM task
LIMIT ? OFFSET ?;