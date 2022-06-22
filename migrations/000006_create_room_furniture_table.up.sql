CREATE TABLE IF NOT EXISTS room_furniture (
    furniture_id BIGSERIAL NOT NULL,
    room_id BIGSERIAL NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    FOREIGN KEY (furniture_id) REFERENCES furniture(furniture_id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES room(room_id) ON DELETE CASCADE,
    PRIMARY KEY (furniture_id, room_id, x, y)
);