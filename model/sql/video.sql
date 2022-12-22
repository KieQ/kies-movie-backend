CREATE TABLE IF NOT EXISTS t_video (
    id serial PRIMARY KEY,
    video_name VARCHAR(255) NOT NULL DEFAULT '',
    video_description TEXT NOT NULL DEFAULT '',
    video_type INTEGER NOT NULL DEFAULT 0,
    region VARCHAR(5) NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    link_type INTEGER NOT NULL DEFAULT 0,
    files TEXT NOT NULL DEFAULT '',
    downloaded BOOLEAN NOT NULL DEFAULT false,
    poster_path VARCHAR(400) NOT NULL DEFAULT '',
    backdrop_path VARCHAR(400) NOT NULL DEFAULT '',
    user_account VARCHAR(255) NOT NULL DEFAULT '',
    tags VARCHAR(200) NOT NULL DEFAULT '',
    liked BOOLEAN NOT NULL DEFAULT false,
    create_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);