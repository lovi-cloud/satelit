package datastore

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A Datastore is type definition of data store.
type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
	GetImage(imageID string) (*europa.BaseImage, error)
	GetImages() ([]europa.BaseImage, error)
	PutImage(image europa.BaseImage) error
	DeleteImage(imageID string) error

	// IPAM
	CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*uuid.UUID, error)
	GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error)
	ListSubnet(ctx context.Context) ([]ipam.Subnet, error)
	DeleteSubnet(ctx context.Context, uuid uuid.UUID) error

	CreateAddress(ctx context.Context, address ipam.Address) (*uuid.UUID, error)
	GetAddressByID(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error)
	ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error)
	DeleteAddress(ctx context.Context, uuid uuid.UUID) error
}
