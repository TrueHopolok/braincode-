INSERT INTO status (owner_name, task_id, score)
VALUES (?,?,?)
ON DUPLICATE KEY UPDATE
score = GREATEST(score, (
	SELECT MAX(s.score)
	FROM status AS s
	WHERE s.owner_name = ?
	AND s.task_id = ?
));