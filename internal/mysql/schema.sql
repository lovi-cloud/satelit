CREATE DATABASE IF NOT EXISTS satelit;

CREATE TABLE IF NOT EXISTS hypervisor (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    iqn VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS volume (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    attached BOOLEAN NOT NULL,
    hostname VARCHAR(255),
    capacity_gb INT NOT NULL,
    base_image_id BINARY(16) NOT NULL,
    host_lun_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS image (
    uuid BINARY(16) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
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
    gateway VARCHAR(255),
    dns_server VARCHAR(255),
    metadata_server VARCHAR(255),
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

CREATE TABLE IF NOT EXISTS lease (
    mac_address VARCHAR(17) NOT NULL PRIMARY KEY,
    address_id BINARY(16) NOT NULL UNIQUE,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    FOREIGN KEY fk_subnet_id(address_id) REFERENCES address(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
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
