SELECT s.id, s.TIMESTAMP, s.task_id, t.title_en, t.title_ru, s.score, COUNT(*) OVER() AS total_amount
FROM Submission AS s
LEFT JOIN Task AS t
ON s.task_id = t.id
WHERE s.owner_name = ?
ORDER BY TIMESTAMP DESC
LIMIT ?;