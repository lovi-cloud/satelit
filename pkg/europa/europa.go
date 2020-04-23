package europa

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type Europa interface {
	CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*Volume, error)
	DeleteVolume(ctx context.Context, name uuid.UUID) error
	ListVolume(ctx context.Context) ([]Volume, error)
	AttachVolume(ctx context.Context, name uuid.UUID, hostname string) (*Volume, error)
}

type Volume struct {
	ID         string
	Attached   bool
	HostName   string
	CapacityGB int
}
