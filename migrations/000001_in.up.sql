CREATE TABLE clients (
    id INT DEFAULT nextval('client_id_seq') PRIMARY KEY,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    encrypt_password VARCHAR NOT NULL
);