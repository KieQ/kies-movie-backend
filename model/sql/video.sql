CREATE TABLE IF NOT EXISTS t_video (
    id serial PRIMARY KEY,
    video_name VARCHAR(255) NOT NULL UNIQUE DEFAULT '',
    video_description TEXT NOT NULL DEFAULT '',
    magnet_link TEXT NOT NULL DEFAULT '',
    video_size INTEGER NOT NULL DEFAULT 0,
    video_type INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(200) NOT NULL DEFAULT '',
    poster_path VARCHAR(200) NOT NULL DEFAULT '',
    backdrop_path VARCHAR(200) NOT NULL DEFAULT '',
    user_account VARCHAR(255) NOT NULL DEFAULT '',
    tag VARCHAR(200) NOT NULL DEFAULT '',
    like_count INTEGER NOT NULL DEFAULT 0,
    create_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);