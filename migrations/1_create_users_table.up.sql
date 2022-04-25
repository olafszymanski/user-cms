CREATE TABLE users (
    id INT PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    email VARCHAR(80) NOT NULL,  
    password VARCHAR(50) NOT NULL,
    admin BIT DEFAULT b'0'
);