CREATE TABLE IF NOT EXISTS transactions(
    id serial PRIMARY KEY,
    from_card_id INT NOT NULL,
    to_card_id INT NOT NULL,
    fee_id INT NULL,
    amount FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);