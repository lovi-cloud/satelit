package europa

import (
	"time"

	uuid "github.com/satori/go.uuid"
	pb "github.com/whywaita/satelit/api/satelit"
)

// A Volume is volume information
type Volume struct {
	ID          string    `db:"id"` // equal HyperMetroPair ID in dorado
	Attached    bool      `db:"attached"`
	HostName    string    `db:"hostname"`
	CapacityGB  uint32    `db:"capacity_gb"`
	BaseImageID uuid.UUID `db:"base_image_id"`
	HostLUNID   int       `db:"host_lun_id"` // set when attached
	BackendName string    `db:"backend_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToPb parse to Satelit API Server pb.Volume
func (v *Volume) ToPb() *pb.Volume {
	pv := &pb.Volume{
		Id:               v.ID,
		Attached:         v.Attached,
		Hostname:         v.HostName,
		CapacityGigabyte: v.CapacityGB,
		BackendName:      v.BackendName,
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
