CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS games (
  id uuid PRIMARY KEY,
  code varchar(16) UNIQUE NOT NULL,
  shift int,
  shift_random boolean DEFAULT false,
  status varchar(20) DEFAULT 'waiting',
  created_at timestamptz DEFAULT now(),
  started_at timestamptz,
  finished_at timestamptz
);

CREATE TABLE IF NOT EXISTS players (
  id uuid PRIMARY KEY,
  game_id uuid REFERENCES games(id) ON DELETE CASCADE,
  name varchar(255) NOT NULL,
  connected boolean DEFAULT false,
  score int DEFAULT 0,
  created_at timestamptz DEFAULT now(),
  finished_at timestamptz
);

CREATE TABLE IF NOT EXISTS questions (
  id uuid PRIMARY KEY,
  game_id uuid REFERENCES games(id) ON DELETE CASCADE,
  plaintext text NOT NULL,
  ciphertext text NOT NULL,
  shift int NOT NULL,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS answers (
  id uuid PRIMARY KEY,
  player_id uuid REFERENCES players(id) ON DELETE CASCADE,
  question_id uuid REFERENCES questions(id) ON DELETE CASCADE,
  is_correct boolean NOT NULL DEFAULT false,
  time_ms int NOT NULL,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS hall_of_fame (
  id uuid PRIMARY KEY,
  player_id uuid,
  game_id uuid,
  player_name varchar(255),
  score int,
  time_taken_ms bigint,
  recorded_at timestamptz DEFAULT now()
);
