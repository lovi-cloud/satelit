CREATE DATABASE IF NOT EXISTS satelit;

USE satelit;

CREATE TABLE IF NOT EXISTS hypervisor (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    iqn VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS image (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    uuid VARCHAR(255) NOT NULL,
    volume_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS subnet (
    uuid BINARY(16) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    network VARCHAR(255) NOT NULL UNIQUE,
    start VARCHAR(255) NOT NULL,
    end VARCHAR(255) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS address (
    uuid BINARY(16) NOT NULL PRIMARY KEY,
    ip VARCHAR(255) NOT NULL UNIQUE,
    subnet_id BINARY(16) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    FOREIGN KEY fk_subnet_id(subnet_id) REFERENCES subnet(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);
