SELECT t.id, t.title, t.owner_name, (
	SELECT MAX(s.score)
	FROM Status AS s
	WHERE t.id = s.task_id
	AND s.owner_name = ?
) AS score, COUNT(*) OVER() AS totalAmount
FROM Task AS t
LIMIT ? OFFSET ?;