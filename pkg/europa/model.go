package europa

import (
	"time"

	uuid "github.com/satori/go.uuid"
	pb "github.com/whywaita/satelit/api/satelit"
)

// A Volume is volume information
type Volume struct {
	ID          string // equal HyperMetroPair ID
	Attached    bool
	HostName    string
	CapacityGB  uint32
	BaseImageID string
	HostLUNID   int       // set when attached
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToPb parse to Satelit API Server pb.Volume
func (v *Volume) ToPb() *pb.Volume {
	pv := &pb.Volume{
		Id:           v.ID,
		Attached:     v.Attached,
		Hostname:     v.HostName,
		CapacityByte: v.CapacityGB,
	}
	return pv
}

// BaseImage is os Image
type BaseImage struct {
	UUID          uuid.UUID `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	CacheVolumeID string    `db:"volume_id"` // if CacheVolumeID is blank, not have cache image
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// ToPb parse to Satelit API Server pb.Image
func (bi *BaseImage) ToPb() *pb.Image {
	i := &pb.Image{
		Id:          bi.UUID.String(),
		Name:        bi.Name,
		VolumeId:    bi.CacheVolumeID,
		Description: bi.Description,
	}

	return i
}
