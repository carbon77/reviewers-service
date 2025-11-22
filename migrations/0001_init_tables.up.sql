CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
  CREATE TYPE pr_status AS ENUM('OPEN', 'MERGED');
EXCEPTION 
  WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users(
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username TEXT NOT NULL UNIQUE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS teams(
  team_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_teams(
  user_id UUID NOT NULL,
  team_id UUID NOT NULL,

  PRIMARY KEY (user_id, team_id),

  FOREIGN KEY (team_id) REFERENCES users (user_id) ON DELETE CASCADE,
  FOREIGN KEY (team_id) REFERENCES teams (team_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pull_requests(
  pull_request_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pull_request_name TEXT NOT NULL,
  status pr_status NOT NULL DEFAULT 'OPEN',
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  merged_at TIMESTAMP,

  author_id UUID NOT NULL,

  FOREIGN KEY (author_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pull_request_reviewers(
  pull_request_id UUID NOT NULL,
  user_id UUID NOT NULL,

  PRIMARY KEY (pull_request_id, user_id),

  FOREIGN KEY (pull_request_id) REFERENCES pull_requests (pull_request_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

