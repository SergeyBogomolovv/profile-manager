CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY NOT NULL,
  email VARCHAR(255) UNIQUE,
  telegram_id BIGINT UNIQUE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE subscription_type AS ENUM ('telegram', 'email');

CREATE TABLE IF NOT EXISTS subscriptions (
  user_id UUID REFERENCES users(id) NOT NULL,
  type subscription_type NOT NULL,
  enabled BOOLEAN DEFAULT TRUE,
  PRIMARY KEY (user_id, type)
);