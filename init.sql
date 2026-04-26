CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    password TEXT,
    role TEXT
);

CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    marca TEXT,
    modelo TEXT,
    ano INT,
    valor NUMERIC
);

CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    car_id INT REFERENCES cars(id) ON DELETE CASCADE,
    ends_at TIMESTAMP NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',
    min_bid NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bids (
    id SERIAL PRIMARY KEY,
    auction_id INT REFERENCES auctions(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id),
    amount NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_logs (
    id SERIAL PRIMARY KEY,
    level TEXT NOT NULL,
    event TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
