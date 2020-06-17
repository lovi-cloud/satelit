package ipam

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// IPAM is interface of IP address management operation.
type IPAM interface {
	CreateSubnet(ctx context.Context, name, prefix, start, end, gateway, dnsServer, metadataServer string) (*Subnet, error)
	GetSubnet(ctx context.Context, uuid uuid.UUID) (*Subnet, error)
	ListSubnet(ctx context.Context) ([]Subnet, error)
	DeleteSubnet(ctx context.Context, uuid uuid.UUID) error

	CreateAddress(ctx context.Context, subnetID uuid.UUID) (*Address, error)
	GetAddress(ctx context.Context, uuid uuid.UUID) (*Address, error)
	ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]Address, error)
	DeleteAddress(ctx context.Context, uuid uuid.UUID) error
}
