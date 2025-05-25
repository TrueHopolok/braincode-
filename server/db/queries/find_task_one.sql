SELECT t.id, t.owner_name, t.title, t.info, (
    SELECT s.score
    FROM Status AS s
    WHERE s.task_id = t.id
    AND s.owner_name = ?
) AS score
FROM Task AS t
WHERE t.id = ?;