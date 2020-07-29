package libvirt

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

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
func (l *Libvirt) CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string) (*ganymede.VirtualMachine, error) {
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
		RootVolumeID:   rootVolumeID,
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

// CreateBridge is
func (l *Libvirt) CreateBridge(ctx context.Context, vlanID uint32) (*ganymede.Bridge, error) {
	clients, err := teleskop.ListClient()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve teleskop clients: %w", err)
	}

	eg := errgroup.Group{}
	for _, client := range clients {
		client := client
		eg.Go(func() error {
			_, err := client.AddVLANInterface(ctx, &agentpb.AddVLANInterfaceRequest{
				VlanId:          vlanID,
				ParentInterface: "bond0",
			})
			if err != nil {
				return err
			}
			_, err = client.AddBridge(ctx, &agentpb.AddBridgeRequest{
				Name: fmt.Sprintf("br%d", vlanID),
			})
			if err != nil {
				return err
			}
			_, err = client.AddInterfaceToBridge(ctx, &agentpb.AddInterfaceToBridgeRequest{
				Bridge:    fmt.Sprintf("br%d", vlanID),
				Interface: fmt.Sprintf("bond0.%d", vlanID),
			})
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	bridge, err := l.ds.CreateBridge(ctx, ganymede.Bridge{
		UUID:   uuid.NewV4(),
		VLANID: vlanID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write bridge: %w", err)
	}

	return bridge, nil
}

// GetBridge is
func (l *Libvirt) GetBridge(ctx context.Context, bridgeID uuid.UUID) (*ganymede.Bridge, error) {
	return l.ds.GetBridge(ctx, bridgeID)
}

// ListBridge is
func (l *Libvirt) ListBridge(ctx context.Context) ([]ganymede.Bridge, error) {
	return l.ds.ListBridge(ctx)
}

// DeleteBridge is
func (l *Libvirt) DeleteBridge(ctx context.Context, bridgeID uuid.UUID) error {
	bridge, err := l.ds.GetBridge(ctx, bridgeID)
	if err != nil {
		return fmt.Errorf("failed to find target bridge: %w", err)
	}

	clients, err := teleskop.ListClient()
	if err != nil {
		return fmt.Errorf("failed to retrieve teleskop clients: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, client := range clients {
		client := client
		eg.Go(func() error {
			_, err := client.DeleteInterfaceFromBridge(ctx, &agentpb.DeleteInterfaceFromBridgeRequest{
				Bridge:    fmt.Sprintf("br%d", bridge.VLANID),
				Interface: fmt.Sprintf("bond0.%d", bridge.VLANID),
			})
			if err != nil {
				return err
			}

			_, err = client.DeleteVLANInterface(ctx, &agentpb.DeleteVLANInterfaceRequest{
				VlanId:          bridge.VLANID,
				ParentInterface: "bond0",
			})
			if err != nil {
				return err
			}

			_, err = client.DeleteBridge(ctx, &agentpb.DeleteBridgeRequest{
				Name: fmt.Sprintf("br%d", bridge.VLANID),
			})
			if err != nil {
				return err
			}

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return l.ds.DeleteBridge(ctx, bridgeID)
}

// AttachInterface is
func (l *Libvirt) AttachInterface(ctx context.Context, vmID, bridgeID, leaseID uuid.UUID, average int, name string) (*ganymede.InterfaceAttachment, error) {
	vm, err := l.ds.GetVirtualMachine(vmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target vm: %w", err)
	}
	bridge, err := l.ds.GetBridge(ctx, bridgeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target bridge: %w", err)
	}
	lease, err := l.ds.GetLeaseByID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target lease: %w", err)
	}

	client, err := teleskop.GetClient(vm.HypervisorName)
	if err != nil {
		return nil, err
	}

	_, err = client.AttachInterface(ctx, &agentpb.AttachInterfaceRequest{
		Uuid:            vm.UUID.String(),
		Bridge:          bridge.Name,
		InboundAverage:  uint32(average),
		OutboundAverage: uint32(average),
		Name:            name,
		MacAddress:      lease.MacAddress.String(),
	})
	if err != nil {
		return nil, err
	}

	attachment, err := l.ds.AttachInterface(ctx, ganymede.InterfaceAttachment{
		UUID:             uuid.NewV4(),
		VirtualMachineID: vm.UUID,
		BridgeID:         bridge.UUID,
		Average:          average,
		Name:             name,
		LeaseID:          lease.UUID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write attachment: %w", err)
	}

	return attachment, nil
}

// DetachInterface is
func (l *Libvirt) DetachInterface(ctx context.Context, attachmentID uuid.UUID) error {
	attachment, err := l.ds.GetAttachment(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("failed to get target attachment: %w", err)
	}
	vm, err := l.ds.GetVirtualMachine(attachment.VirtualMachineID)
	if err != nil {
		return fmt.Errorf("failed to get target vm: %w", err)
	}
	bridge, err := l.ds.GetBridge(ctx, attachment.BridgeID)
	if err != nil {
		return fmt.Errorf("failed to get target bridge: %w", err)
	}
	lease, err := l.ds.GetLeaseByID(ctx, attachment.LeaseID)
	if err != nil {
		return fmt.Errorf("failed to get target lease: %w", err)
	}

	client, err := teleskop.GetClient(vm.HypervisorName)
	if err != nil {
		return err
	}

	_, err = client.DetachInterface(ctx, &agentpb.DetachInterfaceRequest{
		Uuid:            vm.UUID.String(),
		Bridge:          bridge.Name,
		InboundAverage:  uint32(attachment.Average),
		OutboundAverage: uint32(attachment.Average),
		Name:            attachment.Name,
		MacAddress:      lease.MacAddress.String(),
	})
	if err != nil {
		return err
	}

	err = l.ds.DetachInterface(ctx, attachment.UUID)
	if err != nil {
		return err
	}

	return nil
}

// GetAttachment is
func (l *Libvirt) GetAttachment(ctx context.Context, attachmentID uuid.UUID) (*ganymede.InterfaceAttachment, error) {
	return l.ds.GetAttachment(ctx, attachmentID)
}

// ListAttachment is
func (l *Libvirt) ListAttachment(ctx context.Context) ([]ganymede.InterfaceAttachment, error) {
	return l.ds.ListAttachment(ctx)
}
