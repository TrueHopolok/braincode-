SELECT t.id, t.title, t.owner_name, (
	SELECT MAX(s.score)
	FROM Status AS s
	WHERE t.id = s.task_id
	AND s.score = 1
	AND s.owner_name = ?
) AS score, COUNT(*) OVER() AS totalAmount
FROM Task AS t
WHERE t.title LIKE CONCAT(?, '%')
AND (
	t.owner_name = ?
	OR
	true = ?
)
LIMIT ? OFFSET ?;