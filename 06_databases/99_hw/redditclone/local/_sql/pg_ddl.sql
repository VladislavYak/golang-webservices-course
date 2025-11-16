DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(200) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL -- For hashed passwords
);

-- Sample data for testing
INSERT INTO users (login, password) VALUES
('user1', '$2a$10$...'), -- Replace with bcrypt-hashed password
('admin', '$2a$10$...'); -- Replace with bcrypt-hashed password

DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL UNIQUE,
  issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL
);

-- Sample data for testing
INSERT INTO sessions (user_id, token, issued_at, expires_at) VALUES
(1, 'sample_jwt_token_1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour'),
(2, 'sample_jwt_token_2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour');