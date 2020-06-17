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
