CREATE TABLE IF NOT EXISTS cards(
    id serial PRIMARY KEY,
    account_id INT NOT NULL,
    number VARCHAR (16) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);