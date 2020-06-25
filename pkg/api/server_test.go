package api

import (
	"testing"

	"github.com/whywaita/satelit/pkg/ipam/ipam"

	datastoreMemory "github.com/whywaita/satelit/pkg/datastore/memory"
	europaMemory "github.com/whywaita/satelit/pkg/europa/memory"
	ganymedeMemory "github.com/whywaita/satelit/pkg/ganymede/memory"
)

func TestSanitizeImageSize(t *testing.T) {
	inputs := []uint64{
		536870912,  // 0.5 GB
		1073741824, // 1   GB
		1610612736, // 1.5 GB
	}
	outputs := []int{
		1,
		1,
		2,
	}

	for i, input := range inputs {
		o := sanitizeImageSize(input)
		if outputs[i] != o {
			t.Errorf("failed to sanitize (input: %d)", input)
		}
	}
}

// NewMemorySatelit create in-memory Satelit API Server
// for testing Satelit API
func NewMemorySatelit() *SatelitServer {
	europa := europaMemory.New()
	ds := datastoreMemory.New()
	ipamBackend := ipam.New(ds)
	ganymede := ganymedeMemory.New(ds)

	return &SatelitServer{
		Europa:    europa,
		IPAM:      ipamBackend,
		Datastore: ds,
		Ganymede:  ganymede,
	}
}
