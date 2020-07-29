package datastore

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A Datastore is type definition of data store.
type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
	GetImage(imageID uuid.UUID) (*europa.BaseImage, error)
	ListImage() ([]europa.BaseImage, error)
	PutImage(image europa.BaseImage) error
	DeleteImage(imageID uuid.UUID) error

	ListVolume(volumeIDs []string) ([]europa.Volume, error)
	GetVolume(volumeID string) (*europa.Volume, error)
	PutVolume(volume europa.Volume) error
	DeleteVolume(volumeID string) error

	// IPAM
	CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*ipam.Subnet, error)
	GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error)
	ListSubnet(ctx context.Context) ([]ipam.Subnet, error)
	DeleteSubnet(ctx context.Context, uuid uuid.UUID) error

	CreateAddress(ctx context.Context, address ipam.Address) (*ipam.Address, error)
	GetAddressByID(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error)
	ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error)
	DeleteAddress(ctx context.Context, uuid uuid.UUID) error

	CreateLease(ctx context.Context, lease ipam.Lease) (*ipam.Lease, error)
	GetLeaseByID(ctx context.Context, leaseID int) (*ipam.Lease, error)
	GetDHCPLeaseByMACAddress(ctx context.Context, mac types.HardwareAddr) (*ipam.DHCPLease, error)
	ListLease(ctx context.Context) ([]ipam.Lease, error)
	DeleteLease(ctx context.Context, leaseID int) error

	GetVirtualMachine(vmID uuid.UUID) (*ganymede.VirtualMachine, error)
	PutVirtualMachine(vm ganymede.VirtualMachine) error
	DeleteVirtualMachine(vmID uuid.UUID) error

	CreateBridge(ctx context.Context, bridge ganymede.Bridge) (*ganymede.Bridge, error)
	GetBridge(ctx context.Context, bridgeID uuid.UUID) (*ganymede.Bridge, error)
	ListBridge(ctx context.Context) ([]ganymede.Bridge, error)
	DeleteBridge(ctx context.Context, bridgeID uuid.UUID) error

	AttachInterface(ctx context.Context, attachment ganymede.InterfaceAttachment) (*ganymede.InterfaceAttachment, error)
	DetachInterface(ctx context.Context, attachmentID int) error
	GetAttachment(ctx context.Context, attachmentID int) (*ganymede.InterfaceAttachment, error)
	ListAttachment(ctx context.Context) ([]ganymede.InterfaceAttachment, error)
}
