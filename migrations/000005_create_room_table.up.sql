CREATE TABLE IF NOT EXISTS room (
    room_id BIGSERIAL PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    date TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    room_description TEXT NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    room_width INTEGER NOT NULL,
    room_height INTEGER NOT NULL
);