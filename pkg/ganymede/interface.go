package ganymede

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// A Ganymede is type definition of Virtual Machine.
type Ganymede interface {
	CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName string) (*VirtualMachine, error)
}

// VirtualMachine is virtual machine.
type VirtualMachine struct {
	UUID           uuid.UUID `db:"uuid"`
	Name           string    `db:"name"`
	Vcpus          uint32    `db:"vcpus"`
	MemoryKiB      uint64    `db:"memory_kib"`
	HypervisorName string    `db:"hypervisor_name"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}