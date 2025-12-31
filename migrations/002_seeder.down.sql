DELETE FROM wallets
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE email IN ('user1@example.com', 'user2@example.com', 'user3@example.com')
);

DELETE FROM users
WHERE email IN ('user1@example.com', 'user2@example.com', 'user3@example.com');
