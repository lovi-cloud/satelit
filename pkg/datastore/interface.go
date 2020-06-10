package datastore

import (
	"context"

	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
)

// A Datastore is type definition of data store.
type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
	GetImage(imageID string) (*europa.BaseImage, error)
	GetImages() ([]europa.BaseImage, error)
	PutImage(image europa.BaseImage) error
	DeleteImage(imageID string) error
	GetVolume(volumeID string) (*europa.Volume, error)
	PutVolume(volume europa.Volume) error
	DeleteVolume(volumeID string) error
	GetVirtualMachine(vmUUID string) (*ganymede.VirtualMachine, error)
	PutVirtualMachine(vm ganymede.VirtualMachine) error
}
