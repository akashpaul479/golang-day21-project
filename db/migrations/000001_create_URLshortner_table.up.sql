USE urlshortner;
CREATE TABLE urls (
    id VARCHAR(50) PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url VARCHAR(50) NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);