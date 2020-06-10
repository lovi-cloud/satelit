package mysql

import (
	"fmt"

	"github.com/whywaita/satelit/pkg/ganymede"
)

// GetVirtualMachine return virtual machine record
func (m *MySQL) GetVirtualMachine(vmUUID string) (*ganymede.VirtualMachine, error) {
	var vm ganymede.VirtualMachine
	query := fmt.Sprintf(`SELECT 
BIN_TO_UUID(uuid),
name,
vcpus,
memory_kib,
hypervisor_name WHERE uuid = %s`, vmUUID)
	err := m.Conn.Get(&vm, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &vm, nil
}

// PutVirtualMachine write virtual machine record
func (m *MySQL) PutVirtualMachine(vm ganymede.VirtualMachine) error {
	query := `INSERT INTO virtual_machine(name, uuid, vcpus, memory_kib, hypervisor_name) VALUES (?, UUID_TO_BIN(?), ?, ?, ?)`
	_, err := m.Conn.Exec(query, vm.Name, vm.UUID, vm.Vcpus, vm.MemoryKiB, vm.HypervisorName)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// DeleteVirtualMachine delete virtual machine record
func (m *MySQL) DeleteVirtualMachine(vmID string) error {
	return nil
}
