INSERT INTO users (email, username, password, is_verified)
VALUES
  ('user1@example.com', 'user1', '$2a$10$z/R0bv57tue5TaGbqYuEHemghXQYyqCcaPfga4BkkBoRQVOdcZwTy', true),
  ('user2@example.com', 'user2', '$2a$10$z/R0bv57tue5TaGbqYuEHemghXQYyqCcaPfga4BkkBoRQVOdcZwTy', false),
  ('user3@example.com', 'user3', '$2a$10$z/R0bv57tue5TaGbqYuEHemghXQYyqCcaPfga4BkkBoRQVOdcZwTy', true);


INSERT INTO wallets (user_id, balance, currency)
VALUES
  (1, 100000, 'IDR'),
  (2, 250000, 'IDR'),
  (3, 50000, 'IDR');
