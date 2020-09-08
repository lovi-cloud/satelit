package memory

import (
	"context"
	"fmt"
	"time"

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
func (m *Memory) CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string, readBytesSec, writeBytesSec, readIOPSSec, writeIOPSSec uint32, cpuPinningGroupName string) (*ganymede.VirtualMachine, error) {
	u := uuid.NewV4()

	var cpg *ganymede.CPUPinningGroup
	var err error
	if cpuPinningGroupName != "" {
		cpg, err = m.ds.GetCPUPinningGroupByName(ctx, cpuPinningGroupName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve cpu pinning group: %w", err)
		}
	}

	vm := &ganymede.VirtualMachine{
		UUID:           u,
		Name:           name,
		Vcpus:          vcpus,
		MemoryKiB:      memoryKiB,
		HypervisorName: hypervisorName,
		RootVolumeID:   rootVolumeID,
		ReadBytesSec:   readBytesSec,
		WriteBytesSec:  writeBytesSec,
		ReadIOPSSec:    readIOPSSec,
		WriteIOPSSec:   writeIOPSSec,
	}

	if cpg != nil {
		vm.CPUPinningGroupID = cpg.UUID
	}

	err = m.ds.PutVirtualMachine(*vm)
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

// CreateBridge is
func (m *Memory) CreateBridge(ctx context.Context, name string, vlanID uint32) (*ganymede.Bridge, error) {
	return m.ds.CreateBridge(ctx, ganymede.Bridge{
		UUID:      uuid.NewV4(),
		VLANID:    vlanID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// CreateInternalBridge is
func (m *Memory) CreateInternalBridge(ctx context.Context, name string) (*ganymede.Bridge, error) {
	return m.ds.CreateBridge(ctx, ganymede.Bridge{
		UUID:      uuid.NewV4(),
		VLANID:    0,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// GetBridge is
func (m *Memory) GetBridge(ctx context.Context, bridgeID uuid.UUID) (*ganymede.Bridge, error) {
	return m.ds.GetBridge(ctx, bridgeID)
}

// ListBridge is
func (m *Memory) ListBridge(ctx context.Context) ([]ganymede.Bridge, error) {
	return m.ds.ListBridge(ctx)
}

// DeleteBridge is
func (m *Memory) DeleteBridge(ctx context.Context, bridgeID uuid.UUID) error {
	return m.ds.DeleteBridge(ctx, bridgeID)
}

// AttachInterface is
func (m *Memory) AttachInterface(ctx context.Context, vmID, bridgeID, leaseID uuid.UUID, average int, name string) (*ganymede.InterfaceAttachment, error) {
	return m.ds.AttachInterface(ctx, ganymede.InterfaceAttachment{
		UUID:             uuid.NewV4(),
		VirtualMachineID: vmID,
		BridgeID:         bridgeID,
		Average:          average,
		Name:             name,
		LeaseID:          leaseID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	})
}

// DetachInterface is
func (m *Memory) DetachInterface(ctx context.Context, attachmentID uuid.UUID) error {
	return m.ds.DetachInterface(ctx, attachmentID)
}

// GetAttachment is
func (m *Memory) GetAttachment(ctx context.Context, attachmentID uuid.UUID) (*ganymede.InterfaceAttachment, error) {
	return m.ds.GetAttachment(ctx, attachmentID)
}

// ListAttachment is
func (m *Memory) ListAttachment(ctx context.Context) ([]ganymede.InterfaceAttachment, error) {
	return m.ds.ListAttachment(ctx)
}
