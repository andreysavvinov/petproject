-- Users and personal music library
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(5) PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    email VARCHAR(30) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(11) UNIQUE NOT NULL,
    country VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS my_music (
    my_music_id SERIAL PRIMARY KEY,
    users_id VARCHAR(5) UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS albums (
    id VARCHAR(5) PRIMARY KEY,
    album_name VARCHAR(20) NOT NULL,
    release_year VARCHAR(4),
    genre VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS performers (
    id VARCHAR(5) PRIMARY KEY,
    performer_name VARCHAR(30) NOT NULL
);

CREATE TABLE IF NOT EXISTS track_list (
    id VARCHAR(5) PRIMARY KEY,
    track_name VARCHAR(30) NOT NULL,
    genre VARCHAR(20),
    duration VARCHAR(10),
    albums_id VARCHAR(5) REFERENCES albums(id) ON DELETE SET NULL,
    performers_id VARCHAR(5) REFERENCES performers(id) ON DELETE SET NULL,
    my_music_my_music_id INTEGER REFERENCES my_music(my_music_id) ON DELETE SET NULL,
    hidden BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS subscription_validaty (
    users_id VARCHAR(5) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_list_id VARCHAR(5) NOT NULL REFERENCES track_list(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    PRIMARY KEY (users_id, track_list_id)
);
