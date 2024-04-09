-- Active: 1712649823557@@127.0.0.1@5430@rest_api_db@public
CREATE TABLE users (
    id INT PRIMARY KEY not null,
    username VARCHAR(255) not null,
    email VARCHAR(255) not null,
    encrypt_password VARCHAR(255) not null
);