package europa

import (
	"context"

	pb "github.com/whywaita/satelit/api/satelit"

	uuid "github.com/satori/go.uuid"
)

// Europa is interface of volume operation.
type Europa interface {
	CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*Volume, error)
	DeleteVolume(ctx context.Context, id string) error
	ListVolume(ctx context.Context) ([]Volume, error)
	GetVolume(ctx context.Context, id string) (*Volume, error)
	AttachVolume(ctx context.Context, id string, hostname string) error
	DetachVolume(ctx context.Context, id string) error
}

// A Volume is volume information
type Volume struct {
	ID         string
	Attached   bool
	HostName   string
	CapacityGB uint32
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
