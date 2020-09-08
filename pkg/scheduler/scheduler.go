package scheduler

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/ganymede"
)

// Error values for scheduler
var (
	ErrNoValidCoreFound = fmt.Errorf("no valid core found")
	ErrInvalidCorePair  = fmt.Errorf("invalid core pair")
)

// Scheduler is manager of cpu cores.
type Scheduler interface {
	PopCorePair(ctx context.Context, hypervisorID int, numRequestCorePair int, pinningGroupID uuid.UUID) ([]ganymede.CorePair, error)
	PushCorePair(ctx context.Context, pinningGroupID uuid.UUID) error
}
