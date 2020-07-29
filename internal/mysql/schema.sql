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
    base_image_id VARCHAR(36) NOT NULL,
    host_lun_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS image (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    volume_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

CREATE TABLE IF NOT EXISTS subnet (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
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
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    ip VARCHAR(255) NOT NULL UNIQUE,
    subnet_id VARCHAR(36) NOT NULL,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    FOREIGN KEY fk_subnet_id(subnet_id) REFERENCES subnet(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS lease (
    uuid VARCHAR(255) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    mac_address VARCHAR(17) NOT NULL,
    address_id VARCHAR(36) NOT NULL UNIQUE,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    FOREIGN KEY fk_subnet_id(address_id) REFERENCES address(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS virtual_machine (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    vcpus INT NOT NULL,
    memory_kib INT NOT NULL,
    hypervisor_name VARCHAR(255) NOT NULL,
    root_volume_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    FOREIGN KEY fk_root_volume_id(root_volume_id) REFERENCES volume(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS bridge (
    uuid       VARCHAR(36)  NOT NULL PRIMARY KEY,
    vlan_id    INT UNSIGNED NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS interface_attachment (
    uuid VARCHAR(255) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    virtual_machine_id VARCHAR(255) NOT NULL,
    bridge_id VARCHAR(36) NOT NULL,
    average INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    lease_id INT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    FOREIGN KEY fk_subnet_id(lease_id) REFERENCES lease(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY fk_subnet_id(virtual_machine_id) REFERENCES virtual_machine(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);
