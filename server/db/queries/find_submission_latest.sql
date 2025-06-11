SELECT s1.solution
FROM Submission AS s1
WHERE (
    s1.owner_name = ?
    AND s1.task_id = ?
    AND s1.timestamp = ( 
        SELECT MAX(s2.timestamp)
        FROM Submission AS s2
        WHERE s1.owner_name = s2.owner_name AND s2.owner_name = ? AND s2.task_id = ? 
    )
);