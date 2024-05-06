CREATE TABLE clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    encrypt_password VARCHAR NOT NULL
);