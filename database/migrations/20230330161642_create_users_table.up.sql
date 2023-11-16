CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL,
    cellphone VARCHAR (20) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
<<<<<<< HEAD
    updated_at TIMESTAMP NOT NULL
=======
    updated_at TIMESTAMP NOT NULL,

>>>>>>> 5d459e3b2242427f8fb133f65a714d6ae77c113a
);