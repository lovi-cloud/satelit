CREATE DATABASE IF NOT EXISTS satelit;

CREATE TABLE IF NOT EXISTS hypervisor (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    iqn VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    UNIQUE (iqn, hostname),
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

-- Pinning Group
CREATE TABLE IF NOT EXISTS cpu_pinning_group (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    count_of_core INT NOT NULL -- count of request CPU core (can input only a multiple of two)
);

-- all NUMA node had hypervisor_id (not delete when used)
CREATE TABLE IF NOT EXISTS hypervisor_numa_node (
    uuid VARCHAR(36) PRIMARY KEY,
    physical_core_min INT NOT NULL,
    physical_core_max INT NOT NULL,
    logical_core_min INT NOT NULL,
    logical_core_max INT NOT NULL,
    hypervisor_id INT NOT NULL,
    UNIQUE (hypervisor_id, physical_core_min), -- need to unique hypervisor_id and one of four values. physical_core_min don't have specific value in four values.
    FOREIGN KEY fk_hypervisor_id(hypervisor_id) REFERENCES hypervisor(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- all cpu cores had numa_node_id (not delete when used)
CREATE TABLE IF NOT EXISTS hypervisor_cpu_pair (
    uuid VARCHAR(36) PRIMARY KEY,
    numa_node_id VARCHAR(36),
    physical_core_number INT NOT NULL,
    logical_core_number INT NOT NULL,
    UNIQUE (numa_node_id, physical_core_number, logical_core_number),
    FOREIGN KEY fk_numa_node_id(numa_node_id) REFERENCES hypervisor_numa_node(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- already pinned and used CPU core
CREATE TABLE IF NOT EXISTS cpu_core_pinned (
    uuid VARCHAR(36) PRIMARY KEY,
    pinning_group_id VARCHAR(36) NOT NULL,
    hypervisor_cpu_pair_id VARCHAR(36) NOT NULL,
    FOREIGN KEY fk_pinning_group_id(pinning_group_id) REFERENCES cpu_pinning_group(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY fk_cpu_pair_id(hypervisor_cpu_pair_id) REFERENCES hypervisor_cpu_pair(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
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
    vlan_id INT NOT NULL,
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
    uuid VARCHAR(255) NOT NULL PRIMARY KEY,
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
    read_bytes_sec INT UNSIGNED NOT NULL,
    write_bytes_sec INT UNSIGNED NOT NULL,
    read_iops_sec INT UNSIGNED NOT NULL,
    write_iops_sec INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    FOREIGN KEY fk_root_volume_id(root_volume_id) REFERENCES volume(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS bridge (
    uuid       VARCHAR(36)  NOT NULL PRIMARY KEY,
    vlan_id    INT UNSIGNED NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp
);

CREATE TABLE IF NOT EXISTS interface_attachment (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    virtual_machine_id VARCHAR(255) NOT NULL,
    bridge_id VARCHAR(36) NOT NULL,
    average INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    lease_id VARCHAR(36) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    FOREIGN KEY fk_subnet_id(lease_id) REFERENCES lease(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY fk_subnet_id(virtual_machine_id) REFERENCES virtual_machine(uuid) ON DELETE RESTRICT ON UPDATE RESTRICT
);