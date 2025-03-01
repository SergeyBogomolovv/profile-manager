CREATE TABLE IF NOT EXISTS users
(
	user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	email VARCHAR(150) UNIQUE NOT NULL,
	registered_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE account_type AS ENUM ('google', 'credentials');

CREATE TABLE IF NOT EXISTS accounts
(
	user_id UUID REFERENCES users(user_id) NOT NULL,
	provider account_type DEFAULT 'credentials',
	password BYTEA,
	PRIMARY KEY(user_id, provider)
);