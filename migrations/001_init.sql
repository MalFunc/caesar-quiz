-- Skema ini mirror model GORM di internal/models.
-- Aplikasi memakai AutoMigrate, jadi file ini opsional/dokumentasi.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS games (
  id uuid PRIMARY KEY,
  code varchar(16) UNIQUE NOT NULL,
  shift int NOT NULL,
  host_token varchar(64),
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS players (
  id uuid PRIMARY KEY,
  game_id uuid NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  name varchar(255) NOT NULL,
  joined_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS questions (
  id uuid PRIMARY KEY,
  game_id uuid NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  plaintext text NOT NULL,
  cipher text NOT NULL,
  shift int NOT NULL,
  created_at timestamptz DEFAULT now(),
  started_at timestamptz
);

CREATE TABLE IF NOT EXISTS answers (
  id uuid PRIMARY KEY,
  player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  question_id uuid NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  answer text NOT NULL,
  is_correct boolean NOT NULL DEFAULT false,
  answered_at timestamptz DEFAULT now(),
  time_ms bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS hall_of_fames (
  id uuid PRIMARY KEY,
  player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  game_id uuid NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  time_ms bigint NOT NULL DEFAULT 0,
  created_at timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_players_game_id ON players(game_id);
CREATE INDEX IF NOT EXISTS idx_questions_game_id ON questions(game_id);
CREATE INDEX IF NOT EXISTS idx_answers_player_question ON answers(player_id, question_id);
CREATE INDEX IF NOT EXISTS idx_hall_of_fames_game_id ON hall_of_fames(game_id);
