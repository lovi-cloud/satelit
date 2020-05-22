package datastore

import (
	"context"

	"github.com/whywaita/satelit/pkg/europa"
)

// A Datastore is type definition of data store.
type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
	GetImage(imageID string) (*europa.BaseImage, error)
	GetImages() ([]europa.BaseImage, error)
	PutImage(image europa.BaseImage) error
	DeleteImage(imageID string) error
}
