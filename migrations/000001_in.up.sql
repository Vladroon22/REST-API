CREATE TABLE clients (
    id serial PRIMARY KEY,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    encrypt_password VARCHAR NOT NULL
);