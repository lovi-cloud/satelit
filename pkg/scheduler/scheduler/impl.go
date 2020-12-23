package scheduler

import (
	"context"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/pkg/datastore"
	"github.com/lovi-cloud/satelit/pkg/ganymede"
	"github.com/lovi-cloud/satelit/pkg/scheduler"
)

// Scheduler is implementation used datastore.
type Scheduler struct {
	Datastore datastore.Datastore
	mu        sync.RWMutex
}

// New create a new scheduler
func New(ds datastore.Datastore) *Scheduler {
	return &Scheduler{
		Datastore: ds,
	}
}

// PopCorePair allocate pinning cpu core
// numRequestCorePair is number of requested CorePair.
// you can use 2 x numRequestCorePair of cpu cores.
// two is physical core and logical core.
func (s *Scheduler) PopCorePair(ctx context.Context, hypervisorID int, numRequestCorePair int, pinningGroupID uuid.UUID) ([]ganymede.CorePair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodes, err := s.Datastore.GetAvailableCorePair(ctx, hypervisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available core from datastore: %w", err)
	}

	for _, node := range nodes {
		if len(node.CorePairs) < numRequestCorePair {
			// this NUMA node don't have enough cores. go to next node.
			continue
		}

		pairs := getCorePairFromAvailable(node.CorePairs, numRequestCorePair)
		for _, pair := range pairs {
			pinned := ganymede.CPUCorePinned{
				UUID:              uuid.NewV4(),
				CPUPinningGroupID: pinningGroupID,
				CorePairID:        pair.UUID,
			}

			// TODO: bulk insert
			if err := s.Datastore.PutPinnedCore(ctx, pinned); err != nil {
				return nil, fmt.Errorf("failed to pin CPU CorePair: %w", err)
			}
		}

		return pairs, nil
	}

	return nil, scheduler.ErrNoValidCoreFound
}

// PushCorePair free pinned cpu core
func (s *Scheduler) PushCorePair(ctx context.Context, pinningGroupID uuid.UUID) error {
	pinneds, err := s.Datastore.GetPinnedCoreByPinningGroup(ctx, pinningGroupID)
	if err != nil {
		// TODO: check not found, to be skip
		return fmt.Errorf("failed to get pinned core: %w", err)
	}

	for _, pinned := range pinneds {
		// TODO: bulk insert
		if err := s.Datastore.DeletePinnedCore(ctx, pinned.UUID); err != nil {
			return fmt.Errorf("failed to delete pinned core from datastore: %w", err)
		}
	}

	return nil
}

// getCorePairFromAvailable decide pinned cpu cores.
// this function is scheduling logic.
func getCorePairFromAvailable(availableCorePairs []ganymede.CorePair, numRequestCorePair int) []ganymede.CorePair {
	// use head of slice
	return availableCorePairs[:numRequestCorePair]
}
