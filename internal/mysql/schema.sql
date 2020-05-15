CREATE DATABASE IF NOT EXISTS satelit;

USE satelit;

CREATE TABLE IF NOT EXISTS hypervisor (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    iqn VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);