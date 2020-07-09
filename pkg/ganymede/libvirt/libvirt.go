package libvirt

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/ganymede"
	agentpb "github.com/whywaita/teleskop/protoc/agent"
)

// A Libvirt is component of virtual machine.
type Libvirt struct {
	ds datastore.Datastore
}

// New create ganymede application
func New(ds datastore.Datastore) *Libvirt {
	return &Libvirt{
		ds: ds,
	}
}

// CreateVirtualMachine send create operation to teleskop
func (l *Libvirt) CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName string) (*ganymede.VirtualMachine, error) {
	agentReq := &agentpb.AddVirtualMachineRequest{
		Name:       name,
		Vcpus:      vcpus,
		MemoryKib:  memoryKiB,
		BootDevice: bootDeviceName,
	}

	teleskopClient, err := teleskop.GetClient(hypervisorName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve teleskop client: %w", err)
	}
	resp, err := teleskopClient.AddVirtualMachine(ctx, agentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to add virtual machine (name: %s): %w", name, err)
	}

	u, err := uuid.FromString(resp.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response uuid from teleskop: %w", err)
	}

	vm := &ganymede.VirtualMachine{
		Name:           resp.Name,
		UUID:           u,
		Vcpus:          vcpus,
		MemoryKiB:      memoryKiB,
		HypervisorName: hypervisorName,
	}

	err = l.ds.PutVirtualMachine(*vm)
	if err != nil {
		return nil, fmt.Errorf("failed to write virtual machine: %w", err)
	}

	return vm, nil
}

// StartVirtualMachine send start operation to teleskop
func (l *Libvirt) StartVirtualMachine(ctx context.Context, vmID uuid.UUID) error {
	vm, err := l.ds.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to find virtual machine: %w", err)
	}

	teleskopClient, err := teleskop.GetClient(vm.HypervisorName)
	if err != nil {
		return fmt.Errorf("failed to retrieve teleskop client: %w", err)
	}

	_, err = teleskopClient.StartVirtualMachine(ctx, &agentpb.StartVirtualMachineRequest{
		Uuid: vmID.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to start virtual machine: %w", err)
	}

	return nil
}

// DeleteVirtualMachine delete virtual machine
func (l *Libvirt) DeleteVirtualMachine(ctx context.Context, vmID uuid.UUID) error {
	vm, err := l.ds.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to find virtual machine: %w", err)
	}

	// TODO: check alive and boot
	// only stop VM?

	teleskopClient, err := teleskop.GetClient(vm.HypervisorName)
	if err != nil {
		return fmt.Errorf("failed to retrieve teleskop client: %w", err)
	}
	_, err = teleskopClient.DeleteVirtualMachine(ctx, &agentpb.DeleteVirtualMachineRequest{
		Uuid: vmID.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to delete virtual machine from teleskop: %w", err)
	}

	err = l.ds.DeleteVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("failed to delete virtual machine from datastore: %w", err)
	}

	return nil
}
