-- Active: 1712649823557@@127.0.0.1@5430@rest_api_db@public
CREATE TABLE users (
    id AUTO_INCREMENT SERIAL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    encrypt_password VARCHAR(255) NOT NULL
);