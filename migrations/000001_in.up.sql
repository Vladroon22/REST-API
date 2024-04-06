CREATE TABLE users {
    id TINYINT not null unique
    username VARCHAR(255) not null
    email VARCHAR(255) not null
    encrypt_password VARCHAR(255) not null unique
}