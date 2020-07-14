package ganymede

import (
	"context"
	"time"

	pb "github.com/whywaita/satelit/api/satelit"

	uuid "github.com/satori/go.uuid"
)

// A Ganymede is type definition of Virtual Machine.
type Ganymede interface {
	CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string) (*VirtualMachine, error)
	StartVirtualMachine(ctx context.Context, vmID uuid.UUID) error
	DeleteVirtualMachine(ctx context.Context, vmID uuid.UUID) error
}

// VirtualMachine is virtual machine.
type VirtualMachine struct {
	UUID           uuid.UUID `db:"uuid"`
	Name           string    `db:"name"`
	Vcpus          uint32    `db:"vcpus"`
	MemoryKiB      uint64    `db:"memory_kib"`
	HypervisorName string    `db:"hypervisor_name"`
	RootVolumeID   string    `db:"root_volume_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// ToPb convert to type for proto
func (vm *VirtualMachine) ToPb() *pb.VirtualMachine {
	if vm == nil {
		return &pb.VirtualMachine{}
	}

	return &pb.VirtualMachine{
		Uuid:           vm.UUID.String(),
		Name:           vm.Name,
		Vcpus:          vm.Vcpus,
		MemoryKib:      vm.MemoryKiB,
		HypervisorName: vm.HypervisorName,
	}
}
