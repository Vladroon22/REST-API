CREATE TABLE clients (
    id SERIAL,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    encrypt_password VARCHAR NOT NULL
);