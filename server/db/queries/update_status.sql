INSERT INTO Status (owner_name, task_id, score)
VALUES (?,?,?)
ON DUPLICATE KEY UPDATE
score = GREATEST(score, (
	SELECT max_score FROM (
		SELECT MAX(s.score) AS max_score
		FROM Status AS s
		WHERE s.owner_name = ?
		AND s.task_id = ?
	) AS sub
));