package datastore

import (
	"context"
	"net"

	uuid "github.com/satori/go.uuid"

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

	CreateAddress(ctx context.Context, address ipam.Address) (*uuid.UUID, error)
	GetAddressByID(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error)
	ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error)
	DeleteAddress(ctx context.Context, uuid uuid.UUID) error

	CreateLease(ctx context.Context, lease ipam.Lease) (*ipam.Lease, error)
	GetLeaseByMACAddress(ctx context.Context, mac net.HardwareAddr) (*ipam.Lease, error)
	GetDHCPLeaseByMACAddress(ctx context.Context, mac net.HardwareAddr) (*ipam.DHCPLease, error)
	ListLease(ctx context.Context) ([]ipam.Lease, error)
	DeleteLease(ctx context.Context, mac net.HardwareAddr) error

	GetVirtualMachine(vmUUID string) (*ganymede.VirtualMachine, error)
	PutVirtualMachine(vm ganymede.VirtualMachine) error
}
