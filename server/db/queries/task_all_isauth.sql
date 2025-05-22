SELECT t.id, t.title, t.owner_name, EXISTS(
	SELECT *
	FROM status AS s
	WHERE t.id = s.task_id
	AND s.score = 1
	AND s.owner_name = ?
) AS is_solved, COUNT(*) OVER() AS totalAmount
FROM task AS t
LIMIT ? OFFSET ?;