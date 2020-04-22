package europa

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type Europa interface {
	CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*Volume, error)
	ListVolume(ctx context.Context) ([]Volume, error)
}

type Volume struct {
	ID string
}
