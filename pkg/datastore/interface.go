package datastore

import (
	"context"

	"github.com/whywaita/satelit/pkg/europa"
)

// A Datastore is type definition of data store.
type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
	PutImage(image europa.BaseImage) error
}
