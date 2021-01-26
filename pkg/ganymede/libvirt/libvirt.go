package libvirt

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"

	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/internal/client/teleskop"
	"github.com/lovi-cloud/satelit/pkg/datastore"
	"github.com/lovi-cloud/satelit/pkg/ganymede"
	agentpb "github.com/lovi-cloud/teleskop/protoc/agent"
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
func (l *Libvirt) CreateVirtualMachine(ctx context.Context, name string, vcpus uint32, memoryKiB uint64, bootDeviceName, hypervisorName, rootVolumeID string, readBytesSec, writeBytesSec, readIOPSSec, writeIOPSSec uint32, cpuPinningGroupName string) (*ganymede.VirtualMachine, error) {
	agentReq := &agentpb.AddVirtualMachineRequest{
		Name:          name,
		Vcpus:         vcpus,
		MemoryKib:     memoryKiB,
		BootDevice:    bootDeviceName,
		ReadBytesSec:  readBytesSec,
		WriteBytesSec: writeBytesSec,
		ReadIopsSec:   readIOPSSec,
		WriteIopsSec:  writeIOPSSec,
	}

	var cpg *ganymede.CPUPinningGroup
	var err error
	if cpuPinningGroupName != "" {
		cpg, err = l.ds.GetCPUPinningGroupByName(ctx, cpuPinningGroupName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve cpu pinning group: %w", err)
		}

		agentReq.PinningGroupName = cpg.Name
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
		ReadBytesSec:   readBytesSec,
		WriteBytesSec:  writeBytesSec,
		ReadIOPSSec:    readIOPSSec,
		WriteIOPSSec:   writeIOPSSec,
	}

	if cpg != nil {
		vm.CPUPinningGroupID = uuid.NullUUID{
			UUID:  cpg.UUID,
			Valid: true,
		}
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
func (l *Libvirt) CreateBridge(ctx context.Context, name string, vlanID uint32) (*ganymede.Bridge, error) {
	subnet, err := l.ds.GetSubnetByVLAN(ctx, vlanID)
	if err != nil {
		return nil, err
	}
	if subnet.MetadataServer == nil {
		return nil, fmt.Errorf("invalid subnet vlan_id=%d: metadata server must not be null", vlanID)
	}

	mask, _ := subnet.Network.Mask.Size()
	addr := fmt.Sprintf("%s/%d", subnet.MetadataServer.String(), mask)

	clients, err := teleskop.ListClient()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve teleskop clients: %w", err)
	}

	eg := errgroup.Group{}
	for _, client := range clients {
		client := client
		eg.Go(func() error {
			resp, err := client.GetInterfaceName(ctx, &agentpb.GetInterfaceNameRequest{})
			if err != nil {
				return err
			}
			interfaceName := trimVlanID(resp.InterfaceName)

			if _, err := client.AddVLANInterface(ctx, &agentpb.AddVLANInterfaceRequest{
				VlanId:          vlanID,
				ParentInterface: trimVlanID(interfaceName),
			}); err != nil {
				return err
			}

			if _, err = client.AddBridge(ctx, &agentpb.AddBridgeRequest{
				Name:         name,
				MetadataCidr: addr,
				InternalOnly: false,
			}); err != nil {
				return err
			}
			if _, err = client.AddInterfaceToBridge(ctx, &agentpb.AddInterfaceToBridgeRequest{
				Bridge:    name,
				Interface: fmt.Sprintf("%s.%d", interfaceName, vlanID),
			}); err != nil {
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
		Name:   name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write bridge: %w", err)
	}

	return bridge, nil
}

func trimVlanID(interfaceName string) string {
	if !strings.Contains(interfaceName, ".") {
		// interfaceName has not vlan
		return interfaceName
	}

	s := strings.Split(interfaceName, ".")
	return s[0]
}

// CreateInternalBridge is
func (l *Libvirt) CreateInternalBridge(ctx context.Context, name string) (*ganymede.Bridge, error) {
	clients, err := teleskop.ListClient()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve teleskop clients: %w", err)
	}

	eg := errgroup.Group{}
	for _, client := range clients {
		client := client
		eg.Go(func() error {
			_, err = client.AddBridge(ctx, &agentpb.AddBridgeRequest{
				Name:         name,
				MetadataCidr: "",
				InternalOnly: true,
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
		VLANID: 0,
		Name:   name,
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

	var eg errgroup.Group
	for _, client := range clients {
		client := client
		eg.Go(func() error {
			resp, err := client.GetInterfaceName(ctx, &agentpb.GetInterfaceNameRequest{})
			if err != nil {
				return err
			}
			interfaceName := trimVlanID(resp.InterfaceName)

			if bridge.VLANID != 0 {
				_, err := client.DeleteInterfaceFromBridge(ctx, &agentpb.DeleteInterfaceFromBridgeRequest{
					Bridge:    bridge.Name,
					Interface: fmt.Sprintf("%s.%d", interfaceName, bridge.VLANID),
				})
				if err != nil {
					return err
				}

				_, err = client.DeleteVLANInterface(ctx, &agentpb.DeleteVLANInterfaceRequest{
					VlanId:          bridge.VLANID,
					ParentInterface: interfaceName,
				})
				if err != nil {
					return err
				}
			}

			_, err = client.DeleteBridge(ctx, &agentpb.DeleteBridgeRequest{
				Name: bridge.Name,
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
