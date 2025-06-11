-- this query is cursed
SELECT t.id, t.title_en, t.title_ru, t.owner_name, (
	SELECT s.score
	FROM Status AS s
	WHERE t.id = s.task_id
	AND s.owner_name = ?
) AS score, COUNT(*) OVER() AS totalAmount
FROM Task AS t
WHERE (
	CONCAT(t.title_en, t.title_en) LIKE CONCAT('%', ?, '%')
	AND (? OR t.owner_name = ?)
)
LIMIT ? OFFSET ?;