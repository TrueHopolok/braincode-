SELECT t.id, t.title, t.owner_name, 0 AS is_solved, COUNT(*) OVER() AS totalAmount
FROM task AS t
WHERE ? IS NULL OR true
LIMIT ? OFFSET ?;