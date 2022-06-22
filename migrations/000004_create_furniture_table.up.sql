CREATE TABLE IF NOT EXISTS furniture (
    furniture_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(20, 2) NOT NULL,
    furniture_description TEXT NOT NULL DEFAULT '',
    furniture_width INTEGER NOT NULL,
    furniture_height INTEGER NOT NULL,
    image VARCHAR(255) NOT NULL,
    shape INTEGER NOT NULL
);