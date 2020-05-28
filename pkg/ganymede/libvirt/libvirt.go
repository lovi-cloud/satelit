package libvirt

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/ganymede"

	agentpb "github.com/whywaita/satelit/api"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/pkg/datastore"
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
	resp, err := teleskop.GetClient(hypervisorName).AddVirtualMachine(ctx, agentReq)
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
