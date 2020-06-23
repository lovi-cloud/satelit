package ipam

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
)

// A Subnet is subnet information
type Subnet struct {
	UUID           uuid.UUID   `db:"uuid"`
	Name           string      `db:"name"`
	Network        types.IPNet `db:"network"`
	Start          types.IP    `db:"start"`
	End            types.IP    `db:"end"`
	Gateway        *types.IP   `db:"gateway"`
	DNSServer      *types.IP   `db:"dns_server"`
	MetadataServer *types.IP   `db:"metadata_server"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

// A Address is address IP address information
type Address struct {
	UUID      uuid.UUID `db:"uuid"`
	IP        types.IP  `db:"ip"`
	SubnetID  uuid.UUID `db:"subnet_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// A Lease is DHCP lease information
type Lease struct {
	MacAddress types.HardwareAddr `db:"mac_address"`
	AddressID  uuid.UUID          `db:"address_id"`
	CreatedAt  time.Time          `db:"created_at"`
	UpdatedAt  time.Time          `db:"updated_at"`
}

// A DHCPLease is DHCP response information
type DHCPLease struct {
	MacAddress     types.HardwareAddr `db:"mac_address"`
	IP             types.IP           `db:"ip"`
	Network        types.IPNet        `db:"network"`
	Gateway        *types.IP          `db:"gateway"`
	DNSServer      *types.IP          `db:"dns_server"`
	MetadataServer *types.IP          `db:"metadata_server"`
}
