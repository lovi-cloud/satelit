package europa

import (
	"context"

	"github.com/whywaita/satelit/pkg/europa/backend"
)

type Europa struct {
	backend backend.Backend
}

func NewEuropa(b backend.Backend) *Europa {
	return &Europa{b}
}

func (e *Europa) ListVolume(ctx context.Context) ([]Volume, error) {
	return e.backend.ListVolume()
}
