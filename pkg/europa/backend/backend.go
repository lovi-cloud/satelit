package backend

import "github.com/whywaita/satelit/pkg/europa"

type Backend interface {
	CreateVolume(capacity int) (*europa.Volume, error)
	ListVolume() ([]europa.Volume, error)
}
