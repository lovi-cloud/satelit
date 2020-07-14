package memory

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/ganymede"
)

// Memory is backend of ganymede by in-memory for testing
type Memory struct {
	ds datastore.Datastore
}

// New create ganymade application
func New(ds datastore.Datastore) *Memory {
	return &Memory{ds: ds}
}

// CreateVirtualMachine add virtual machine
func (m *Memory) CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string) (*ganymede.VirtualMachine, error) {
	u := uuid.NewV4()

	vm := &ganymede.VirtualMachine{
		UUID:           u,
		Name:           name,
		Vcpus:          vcpus,
		MemoryKiB:      memoryKiB,
		HypervisorName: hypervisorName,
		RootVolumeID:   rootVolumeID,
	}

	err := m.ds.PutVirtualMachine(*vm)
	if err != nil {
		return nil, fmt.Errorf("failed to write virtual machine: %w", err)
	}

	return vm, nil
}

// StartVirtualMachine start virtual machine
func (m *Memory) StartVirtualMachine(ctx context.Context, vmID uuid.UUID) error {
	_, err := m.ds.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to find virtual machine: %w", err)
	}

	return nil
}

// DeleteVirtualMachine delete virtual machine
func (m *Memory) DeleteVirtualMachine(ctx context.Context, vmID uuid.UUID) error {
	_, err := m.ds.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to find virtual machine: %w", err)
	}

	err = m.ds.DeleteVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to delete virtual machine from datastore: %w", err)
	}

	return nil
}
