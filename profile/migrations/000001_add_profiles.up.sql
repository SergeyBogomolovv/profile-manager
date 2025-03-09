CREATE TYPE user_gender AS ENUM ('male', 'female', 'not specified');

CREATE TABLE IF NOT EXISTS profiles (
  user_id BIGINT PRIMARY KEY NOT NULL,
  username VARCHAR(255) NOT NULL UNIQUE,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  birth_date DATE,
  gender user_gender DEFAULT 'not specified',
  avatar TEXT
);