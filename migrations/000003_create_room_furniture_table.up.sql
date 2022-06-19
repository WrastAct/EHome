CREATE TABLE IF NOT EXISTS room_furniture (
    room_furniture_id BIGSERIAL PRIMARY KEY,
    furniture_id BIGSERIAL NOT NULL,
    room_id BIGSERIAL NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    FOREIGN KEY (furniture_id) REFERENCES furniture(furniture_id),
    FOREIGN KEY (room_id) REFERENCES room(room_id)
);