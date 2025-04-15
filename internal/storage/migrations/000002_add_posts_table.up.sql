CREATE TABLE IF NOT EXISTS posts(
  postid SERIAL PRIMARY KEY,
  userid INT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
  post_content TEXT NOT NULL,
  image_url TEXT,
  created_at TIMESTAMP DEFAULT now()
);
