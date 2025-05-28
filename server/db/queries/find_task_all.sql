SELECT t.id, t.title_en, t.title_ru, t.owner_name, (
	SELECT s.score
	FROM Status AS s
	WHERE t.id = s.task_id
	AND s.owner_name = ?
) AS score, COUNT(*) OVER() AS totalAmount
FROM Task AS t
WHERE (
	t.title_en LIKE CONCAT(?, '%')
	OR
	t.title_ru LIKE CONCAT(?, '%')
)
AND (
	t.owner_name = ?
	OR
	true = ?
)
LIMIT ? OFFSET ?;