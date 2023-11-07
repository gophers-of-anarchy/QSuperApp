CREATE TABLE IF NOT EXISTS accounts(
    id serial PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR (50) NOT NULL,
    balance FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);