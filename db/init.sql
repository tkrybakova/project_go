CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    slot_id VARCHAR(255),
    date VARCHAR(50),
    status VARCHAR(50)
);
CREATE TABLE brigades (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    brigade_id INT REFERENCES brigades(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);
