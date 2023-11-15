CREATE TABLE IF NOT EXISTS transaction_fees(
   id serial PRIMARY KEY,
   transaction_id INT NOT NULL,
   amount FLOAT NOT NULL,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL
);