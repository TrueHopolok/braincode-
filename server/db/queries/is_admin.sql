SELECT EXISTS(
    SELECT *
    FROM User
    WHERE name=? AND is_admin
)