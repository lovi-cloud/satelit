CREATE DATABASE IF NOT EXISTS satelit;

CREATE TABLE IF NOT EXISTS hypervisor (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    iqn VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS image (
    uuid BINARY(16) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    volume_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS virtual_machine (
    uuid BINARY(16) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    vcpus INT NOT NULL,
    memory_kib INT NOT NULL,
    hypervisor_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
)