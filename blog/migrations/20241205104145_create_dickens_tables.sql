-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY,
   user_id INTEGER,
   name TEXT NOT NULL,
   publish BOOLEAN,
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   FOREIGN KEY (user_id) 
      REFERENCES users (id) 
         ON DELETE CASCADE 
         ON UPDATE NO ACTION
    
);

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
   user_id INTEGER,
   body TEXT NOT NULL,
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   FOREIGN KEY (user_id) 
      REFERENCES users (id) 
         ON DELETE CASCADE 
         ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS blogs (
	id INTEGER PRIMARY KEY,
   user_id INTEGER,
   title TEXT NOT NULL,
   body TEXT NOT NULL,
   publish BOOLEAN,
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   FOREIGN KEY (user_id) 
      REFERENCES users (id) 
         ON DELETE CASCADE 
         ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS category_blogs (
   category_id INTEGER NOT NULL,
   blog_id INTEGER NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (category_id, blog_id),
   FOREIGN KEY (category_id) 
      REFERENCES categories (id) 
        ON DELETE CASCADE 
        ON UPDATE NO ACTION,
   FOREIGN KEY (blog_id) 
      REFERENCES blogs (id) 
        ON DELETE CASCADE 
        ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS comment_blogs (
   comment_id INTEGER NOT NULL,
   blog_id INTEGER NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (comment_id, blog_id),
   FOREIGN KEY (comment_id) 
      REFERENCES comments (id) 
        ON DELETE CASCADE 
        ON UPDATE NO ACTION,
   FOREIGN KEY (blog_id) 
      REFERENCES blogs (id) 
        ON DELETE CASCADE 
        ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS profiles(
	id INTEGER PRIMARY KEY,
   user_id INTEGER,
   username TEXT NOT NULL,
	image TEXT,
	bio TEXT,
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   FOREIGN KEY (user_id) 
      REFERENCES users (id) 
         ON DELETE CASCADE 
         ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comment_blogs;
DROP TABLE IF EXISTS category_blogs;
DROP TABLE IF EXISTS blogs;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS profiles;
-- +goose StatementEnd
