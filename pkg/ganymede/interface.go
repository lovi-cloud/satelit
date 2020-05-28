package ganymede

import "context"

// A Ganymede is type definition of Virtual Machine.
type Ganymede interface {
	CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName string) (*VirtualMachine, error)
}

// VirtualMachine is virtual machine.
type VirtualMachine struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	UUID           string `db:"uuid"`
	Vcpus          uint32 `db:"vcpus"`
	MemoryKiB      uint64 `db:"memory_kib"`
	HypervisorName string `db:"hypervisor_name"`
}
