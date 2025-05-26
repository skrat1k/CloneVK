CREATE TABLE IF NOT EXISTS follows(
  ID SERIAL PRIMARY KEY,
  follower_id INT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
  followed_id INT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT now(),
  UNIQUE (follower_id, followed_id),
  CHECK (follower_id <> followed_id)
)
