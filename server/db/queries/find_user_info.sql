/*
Acceptance rate 	= (COUNT(submission WHERE score = 1) / COUNT(submission))
Solved rate 		= (COUNT(status WHERE score = 1) / COUNT(task))
*/
SELECT
	(
		SUM(
			CASE WHEN s.score = 1 
				THEN 1 
				ELSE 0 
			END
		) / COUNT(s.id)
	) AS acceptance_rate,
	(
		(
			SELECT COUNT(*) 
			FROM status AS st 
			WHERE st.score = 1
		) * 1.0 / (
			SELECT COUNT(*) 
			FROM task
		)
	) AS solved_rate
FROM submission AS s
WHERE s.owner_name = ?;