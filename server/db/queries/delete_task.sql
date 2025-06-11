DELETE FROM Task AS t
WHERE (
    t.owner_name = ?
    OR
    EXISTS (
        SELECT *
        FROM User AS u
        WHERE u.name = ?
        AND u.is_admin
    )
) AND t.id = ?;