package europa

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// Europa is interface of volume operation.
type Europa interface {
	CreateVolumeRaw(ctx context.Context, name uuid.UUID, capacity int) (*Volume, error)
	CreateVolumeImage(ctx context.Context, name uuid.UUID, capacity int, imageID string) (*Volume, error)
	DeleteVolume(ctx context.Context, id string) error
	ListVolume(ctx context.Context) ([]Volume, error)
	GetVolume(ctx context.Context, id string) (*Volume, error)
	AttachVolume(ctx context.Context, id string, hostname string, isTeleskop bool) (int, string, error)
	DetachVolume(ctx context.Context, id string) error
	UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*BaseImage, error)
	DeleteImage(ctx context.Context, id string) error
}
