package ganymede

import (
	"context"
	"time"

	pb "github.com/whywaita/satelit/api/satelit"
	dspb "github.com/whywaita/satelit/api/satelit_datastore"

	uuid "github.com/satori/go.uuid"
)

// A Ganymede is type definition of Virtual Machine.
type Ganymede interface {
	CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string, readBytesSec, writeBytesSec, readIOPSSec, writeIOPSSec uint32, cpuPinningGroupName string) (*VirtualMachine, error)
	StartVirtualMachine(ctx context.Context, vmID uuid.UUID) error
	DeleteVirtualMachine(ctx context.Context, vmID uuid.UUID) error

	CreateBridge(ctx context.Context, name string, vlanID uint32) (*Bridge, error)
	CreateInternalBridge(ctx context.Context, name string) (*Bridge, error)
	GetBridge(ctx context.Context, bridgeID uuid.UUID) (*Bridge, error)
	ListBridge(ctx context.Context) ([]Bridge, error)
	DeleteBridge(ctx context.Context, bridgeID uuid.UUID) error

	AttachInterface(ctx context.Context, vmID, bridgeID, leaseID uuid.UUID, average int, name string) (*InterfaceAttachment, error)
	DetachInterface(ctx context.Context, attachmentID uuid.UUID) error
	GetAttachment(ctx context.Context, attachmentID uuid.UUID) (*InterfaceAttachment, error)
	ListAttachment(ctx context.Context) ([]InterfaceAttachment, error)
}

// VirtualMachine is virtual machine.
type VirtualMachine struct {
	UUID              uuid.UUID `db:"uuid"`
	Name              string    `db:"name"`
	Vcpus             uint32    `db:"vcpus"`
	MemoryKiB         uint64    `db:"memory_kib"`
	HypervisorName    string    `db:"hypervisor_name"`
	RootVolumeID      string    `db:"root_volume_id"`
	RootVolumeGB      uint32    `db:"capacity_gb"`
	ReadBytesSec      uint32    `db:"read_bytes_sec"`
	WriteBytesSec     uint32    `db:"write_bytes_sec"`
	ReadIOPSSec       uint32    `db:"read_iops_sec"`
	WriteIOPSSec      uint32    `db:"write_iops_sec"`
	SourceImageID     uuid.UUID `db:"base_image_id"`
	CPUPinningGroupID uuid.UUID `db:"cpu_pinning_group_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// ToPb convert to type for proto
func (vm *VirtualMachine) ToPb(cpuPinningGroupName string) *pb.VirtualMachine {
	if vm == nil {
		return &pb.VirtualMachine{}
	}

	return &pb.VirtualMachine{
		Uuid:             vm.UUID.String(),
		Name:             vm.Name,
		Vcpus:            vm.Vcpus,
		MemoryKib:        vm.MemoryKiB,
		HypervisorName:   vm.HypervisorName,
		RootVolumeGb:     vm.RootVolumeGB,
		ReadBytesSec:     vm.ReadBytesSec,
		WriteBytesSec:    vm.WriteBytesSec,
		ReadIopsSec:      vm.ReadIOPSSec,
		WriteIopsSec:     vm.WriteIOPSSec,
		SourceImageId:    vm.SourceImageID.String(),
		PinningGroupName: cpuPinningGroupName,
	}
}

// Bridge is bridge.
type Bridge struct {
	UUID      uuid.UUID `db:"uuid"`
	VLANID    uint32    `db:"vlan_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToPb is
func (b *Bridge) ToPb() *pb.Bridge {
	if b == nil {
		return &pb.Bridge{}
	}

	return &pb.Bridge{
		Uuid:   b.UUID.String(),
		VlanId: b.VLANID,
		Name:   b.Name,
	}
}

// InterfaceAttachment is
type InterfaceAttachment struct {
	UUID             uuid.UUID `db:"uuid"`
	VirtualMachineID uuid.UUID `db:"virtual_machine_id"`
	BridgeID         uuid.UUID `db:"bridge_id"`
	Average          int       `db:"average"`
	Name             string    `db:"name"`
	LeaseID          uuid.UUID `db:"lease_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// ToPb is
func (i *InterfaceAttachment) ToPb() *pb.InterfaceAttachment {
	if i == nil {
		return &pb.InterfaceAttachment{}
	}

	return &pb.InterfaceAttachment{
		Uuid:             i.UUID.String(),
		VirtualMachineId: i.VirtualMachineID.String(),
		BridgeId:         i.BridgeID.String(),
		Average:          int64(i.Average),
		Name:             i.Name,
		LeaseId:          i.LeaseID.String(),
	}
}

// HyperVisor is host of virtual machines
type HyperVisor struct {
	ID        int       `db:"id"`
	IQN       string    `db:"iqn"`
	Hostname  string    `db:"hostname"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// NUMANode is Non-uniform memory access in hypervisor
type NUMANode struct {
	UUID            uuid.UUID `db:"uuid"`
	CorePairs       []CorePair
	PhysicalCoreMin uint32    `db:"physical_core_min"`
	PhysicalCoreMax uint32    `db:"physical_core_max"`
	LogicalCoreMin  uint32    `db:"logical_core_min"`
	LogicalCoreMax  uint32    `db:"logical_core_max"`
	HypervisorID    int       `db:"hypervisor_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// ToPb convert to type for proto
func (node *NUMANode) ToPb() *dspb.NumaNode {
	var pairs []*dspb.CorePair
	for _, p := range node.CorePairs {
		pairs = append(pairs, p.ToPb())
	}

	return &dspb.NumaNode{
		Pairs:           pairs,
		PhysicalCoreMin: node.PhysicalCoreMin,
		PhysicalCoreMax: node.PhysicalCoreMax,
		LogicalCoreMin:  node.LogicalCoreMin,
		LogicalCoreMax:  node.LogicalCoreMax,
	}
}

// CorePair is pair of cpu core
type CorePair struct {
	UUID         uuid.UUID `db:"uuid"`
	PhysicalCore uint32    `db:"physical_core_number"`
	LogicalCore  uint32    `db:"logical_core_number"`
	NUMANodeID   uuid.UUID `db:"numa_node_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToPb convert to type for proto
func (cp *CorePair) ToPb() *dspb.CorePair {
	return &dspb.CorePair{
		PhysicalCore: cp.PhysicalCore,
		LogicalCore:  cp.LogicalCore,
	}
}

// CPUPinningGroup is group of cpu cores.
// use a same group's cpu core if joined a same group.
type CPUPinningGroup struct {
	UUID         uuid.UUID `db:"uuid"`
	Name         string    `db:"name"`
	HypervisorID int       `db:"hypervisor_id"`
	CountCore    int       `db:"count_of_core"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ToPb convert to type for proto
func (cpg *CPUPinningGroup) ToPb() *pb.CPUPinningGroup {
	return &pb.CPUPinningGroup{
		Uuid:        cpg.UUID.String(),
		Name:        cpg.Name,
		CountOfCore: uint32(cpg.CountCore),
	}
}

// CPUCorePinned is pinned cpu
type CPUCorePinned struct {
	UUID              uuid.UUID `db:"uuid"`
	CPUPinningGroupID uuid.UUID `db:"pinning_group_id"`
	CorePairID        uuid.UUID `db:"hypervisor_cpu_pair_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
