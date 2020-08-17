package mysql

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/ganymede"
)

// GetVirtualMachine return virtual machine record
func (m *MySQL) GetVirtualMachine(vmID uuid.UUID) (*ganymede.VirtualMachine, error) {
	var vm ganymede.VirtualMachine
	query := fmt.Sprintf(`SELECT uuid, name, vcpus, memory_kib, hypervisor_name, root_volume_id, volume.capacity_gb, read_bytes_sec, write_bytes_sec, read_iops_sec, write_iops_sec, volume.base_image_id FROM virtual_machine JOIN volume ON virtual_machine.root_volume_id = volume.id WHERE virtual_machine.uuid = '%s'`, vmID.String())
	err := m.Conn.Get(&vm, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &vm, nil
}

// PutVirtualMachine write virtual machine record
func (m *MySQL) PutVirtualMachine(vm ganymede.VirtualMachine) error {
	query := `INSERT INTO virtual_machine(name, uuid, vcpus, memory_kib, hypervisor_name, root_volume_id, read_bytes_sec, write_bytes_sec, read_iops_sec, write_iops_sec) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.Conn.Exec(query, vm.Name, vm.UUID, vm.Vcpus, vm.MemoryKiB, vm.HypervisorName, vm.RootVolumeID, vm.ReadBytesSec, vm.WriteBytesSec, vm.ReadIOPSSec, vm.WriteIOPSSec)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// DeleteVirtualMachine delete virtual machine record
func (m *MySQL) DeleteVirtualMachine(vmID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM virtual_machine WHERE uuid = "%s"`, vmID.String())
	_, err := m.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}

// GetHostnameByAddress is
func (m *MySQL) GetHostnameByAddress(address types.IP) (string, error) {
	query := `SELECT v.name FROM address AS a JOIN subnet AS s ON a.subnet_id = s.uuid JOIN lease AS l ON a.uuid = l.address_id JOIN interface_attachment AS i ON l.uuid = i.lease_id JOIN virtual_machine AS v ON i.virtual_machine_id = v.uuid WHERE ip = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}
	var hostname string
	err = stmt.Get(&hostname, address)
	if err != nil {
		return "", fmt.Errorf("failed to get hostname by address: %w", err)
	}

	return hostname, nil
}
