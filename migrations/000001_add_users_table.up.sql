CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email varchar(255) NOT NULL UNIQUE,
  nickname varchar(255) NOT NULL UNIQUE,
  password varchar(512) NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX users_id_index ON users (id);
CREATE INDEX users_nickname_index ON users (nickname);
