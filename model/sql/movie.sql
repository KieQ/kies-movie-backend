CREATE TABLE IF NOT EXISTS t_movie (
    id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE DEFAULT '',
    magnet_link TEXT NOT NULL DEFAULT '',
    size INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(200) NOT NULL DEFAULT '',
    profile VARCHAR(200) NOT NULL DEFAULT '',
    user_account VARCHAR(255) NOT NULL DEFAULT '',
    is_private BOOLEAN NOT NULL DEFAULT true,
    tag VARCHAR(200) NOT NULL DEFAULT '',
    like_count INTEGER NOT NULL DEFAULT 0,
    create_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);