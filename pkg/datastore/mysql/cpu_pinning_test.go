package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/jmoiron/sqlx"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/internal/testutils"
	"github.com/whywaita/satelit/pkg/ganymede"
)

const (
	testCPUPinningGroupUUID = "0c817251-5f39-4983-b064-9658ecaf5191"
	testCPUPinningGroupName = "testcpg"
	testCPUPinnedCoreUUID   = "144a38e6-4aa5-41d7-8581-2036f2312da2"
	testNUMANodeUUID        = "162b42f5-2eea-4fd1-b57b-c598db69fb4a"
	testCorePairUUID        = "bbd07e10-75c5-4ed4-986b-500320164310"
)

func TestMySQL_PutCPUPinningGroup(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}

	tests := []struct {
		input ganymede.CPUPinningGroup
		want  *ganymede.CPUPinningGroup
		err   bool
	}{
		{
			input: ganymede.CPUPinningGroup{
				UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
				Name:         testCPUPinningGroupName,
				HypervisorID: hypervisorID,
				CountCore:    4,
			},
			want: &ganymede.CPUPinningGroup{
				UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
				Name:         testCPUPinningGroupName,
				HypervisorID: hypervisorID,
				CountCore:    4,
			},
			err: false,
		},
	}

	for _, test := range tests {
		err := testDatastore.PutCPUPinningGroup(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		got, err := getCPUPinningGroupFromSQL(testDB, test.input.UUID)
		if err != nil {
			t.Fatalf("failed to get cpu pinning group from SQL: %+v", err)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_GetCPUPinningGroup(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	cpg := ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}
	err = testDatastore.PutCPUPinningGroup(context.Background(), cpg)
	if err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.CPUPinningGroup
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testCPUPinningGroupUUID),
			want: &ganymede.CPUPinningGroup{
				UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
				Name:         testCPUPinningGroupName,
				HypervisorID: hypervisorID,
				CountCore:    4,
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetCPUPinningGroup(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_GetCPUPinningGroupByName(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	cpg := ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}
	err = testDatastore.PutCPUPinningGroup(context.Background(), cpg)
	if err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}

	tests := []struct {
		input string
		want  *ganymede.CPUPinningGroup
		err   bool
	}{
		{
			input: testCPUPinningGroupName,
			want: &ganymede.CPUPinningGroup{
				UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
				Name:         testCPUPinningGroupName,
				HypervisorID: hypervisorID,
				CountCore:    4,
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetCPUPinningGroupByName(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_DeleteCPUPinningGroup(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	cpg := ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}
	err = testDatastore.PutCPUPinningGroup(context.Background(), cpg)
	if err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.CPUPinningGroup
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testCPUPinningGroupUUID),
			want:  nil,
			err:   false,
		},
	}

	for _, test := range tests {
		err := testDatastore.DeleteCPUPinningGroup(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}

		got, err := getCPUPinningGroupFromSQL(testDB, uuid.FromStringOrNil(testCPUPinningGroupUUID))
		if test.want != got {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestMySQL_GetAvailableCorePair(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	testCorePairUUIDs := []uuid.UUID{
		uuid.FromStringOrNil("9cf11645-ec85-4607-b638-cd592819bbae"),
		uuid.FromStringOrNil("25b403a9-cdd7-4176-8d44-c922220bdcb8"),
		uuid.FromStringOrNil("2cc61359-8912-4187-aadc-8692574b1b52"),
		uuid.FromStringOrNil("e77523a3-fef0-4864-b24f-4f9579a65eed"),
	}

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	if err := testDatastore.PutCPUPinningGroup(context.Background(), ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    2,
	}); err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}
	nodes := []ganymede.NUMANode{
		ganymede.NUMANode{
			UUID: uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs: []ganymede.CorePair{
				ganymede.CorePair{
					UUID:         testCorePairUUIDs[0],
					PhysicalCore: 0,
					LogicalCore:  5,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
				ganymede.CorePair{
					UUID:         testCorePairUUIDs[1],
					PhysicalCore: 1,
					LogicalCore:  6,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
				ganymede.CorePair{
					UUID:         testCorePairUUIDs[2],
					PhysicalCore: 2,
					LogicalCore:  7,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
				ganymede.CorePair{
					UUID:         testCorePairUUIDs[3],
					PhysicalCore: 3,
					LogicalCore:  8,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
			},
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 3,
			LogicalCoreMin:  5,
			LogicalCoreMax:  8,
			HypervisorID:    hypervisorID,
		},
	}
	if err = testDatastore.PutHypervisorNUMANode(context.Background(), nodes, hypervisorID); err != nil {
		t.Fatalf("failed to put hypervisor numa node: %+v", err)
	}

	tests := []struct {
		input []uuid.UUID // pinned CorePairUUID
		want  []ganymede.NUMANode
		err   bool
	}{
		{
			input: nil,
			want: []ganymede.NUMANode{
				ganymede.NUMANode{
					UUID: uuid.FromStringOrNil(testNUMANodeUUID),
					CorePairs: []ganymede.CorePair{
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[0],
							PhysicalCore: 0,
							LogicalCore:  5,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[1],
							PhysicalCore: 1,
							LogicalCore:  6,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[2],
							PhysicalCore: 2,
							LogicalCore:  7,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[3],
							PhysicalCore: 3,
							LogicalCore:  8,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
					},
				},
			},
			err: false,
		},
		{
			input: []uuid.UUID{
				testCorePairUUIDs[0],
				testCorePairUUIDs[1],
			},
			want: []ganymede.NUMANode{
				ganymede.NUMANode{
					UUID: uuid.FromStringOrNil(testNUMANodeUUID),
					CorePairs: []ganymede.CorePair{
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[2],
							PhysicalCore: 2,
							LogicalCore:  7,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
						ganymede.CorePair{
							UUID:         testCorePairUUIDs[3],
							PhysicalCore: 3,
							LogicalCore:  8,
							NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
						},
					},
				},
			},
			err: false,
		},
		{
			input: []uuid.UUID{
				testCorePairUUIDs[2],
				testCorePairUUIDs[3],
			},
			want: []ganymede.NUMANode{
				ganymede.NUMANode{
					UUID:      uuid.FromStringOrNil(testNUMANodeUUID),
					CorePairs: nil,
				},
			},
			err: false,
		},
	}

	for _, test := range tests {
		if len(test.input) > 0 {
			for _, i := range test.input {
				// pinning
				if err := testDatastore.PutPinnedCore(context.Background(), ganymede.CPUCorePinned{
					UUID:              uuid.NewV4(),
					CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
					CorePairID:        i,
				}); err != nil {
					t.Fatalf("failed to put pinned core: %+v", err)
				}
			}
		}

		got, err := testDatastore.GetAvailableCorePair(context.Background(), hypervisorID)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got = setFakeTimeNUMANodes(got)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_GetCPUCorePair(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	nodes := []ganymede.NUMANode{
		ganymede.NUMANode{
			UUID: uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs: []ganymede.CorePair{
				ganymede.CorePair{
					UUID:         uuid.FromStringOrNil(testCorePairUUID),
					PhysicalCore: 0,
					LogicalCore:  41,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
			},
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 40,
			LogicalCoreMin:  41,
			LogicalCoreMax:  80,
			HypervisorID:    hypervisorID,
		},
	}
	err = testDatastore.PutHypervisorNUMANode(context.Background(), nodes, hypervisorID)
	if err != nil {
		t.Fatalf("failed to put hypervisor numa node: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.CorePair
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testCorePairUUID),
			want: &ganymede.CorePair{
				UUID:         uuid.FromStringOrNil(testCorePairUUID),
				PhysicalCore: 0,
				LogicalCore:  41,
				NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetCPUCorePair(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_GetPinnedCoreByPinningGroup(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	if err := testDatastore.PutCPUPinningGroup(context.Background(), ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}); err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}
	nodes := []ganymede.NUMANode{
		ganymede.NUMANode{
			UUID: uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs: []ganymede.CorePair{
				ganymede.CorePair{
					UUID:         uuid.FromStringOrNil(testCorePairUUID),
					PhysicalCore: 0,
					LogicalCore:  41,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
			},
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 40,
			LogicalCoreMin:  41,
			LogicalCoreMax:  80,
			HypervisorID:    hypervisorID,
		},
	}
	if err = testDatastore.PutHypervisorNUMANode(context.Background(), nodes, hypervisorID); err != nil {
		t.Fatalf("failed to put hypervisor numa node: %+v", err)
	}

	pinned := ganymede.CPUCorePinned{
		UUID:              uuid.FromStringOrNil(testCPUPinnedCoreUUID),
		CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
		CorePairID:        uuid.FromStringOrNil(testCorePairUUID),
	}
	if err := testDatastore.PutPinnedCore(context.Background(), pinned); err != nil {
		t.Fatalf("failed to put pinned core: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  []ganymede.CPUCorePinned
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testCPUPinningGroupUUID),
			want: []ganymede.CPUCorePinned{
				ganymede.CPUCorePinned{
					UUID:              uuid.FromStringOrNil(testCPUPinnedCoreUUID),
					CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
					CorePairID:        uuid.FromStringOrNil(testCorePairUUID),
				},
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetPinnedCoreByPinningGroup(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got[0].CreatedAt = time.Time{}
			got[0].UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_PutPinnedCore(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	if err := testDatastore.PutCPUPinningGroup(context.Background(), ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}); err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}
	nodes := []ganymede.NUMANode{
		ganymede.NUMANode{
			UUID: uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs: []ganymede.CorePair{
				ganymede.CorePair{
					UUID:         uuid.FromStringOrNil(testCorePairUUID),
					PhysicalCore: 0,
					LogicalCore:  41,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
			},
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 40,
			LogicalCoreMin:  41,
			LogicalCoreMax:  80,
			HypervisorID:    hypervisorID,
		},
	}
	if err = testDatastore.PutHypervisorNUMANode(context.Background(), nodes, hypervisorID); err != nil {
		t.Fatalf("failed to put hypervisor numa node: %+v", err)
	}

	tests := []struct {
		input ganymede.CPUCorePinned
		want  *ganymede.CPUCorePinned
		err   bool
	}{
		{
			input: ganymede.CPUCorePinned{
				UUID:              uuid.FromStringOrNil(testCPUPinnedCoreUUID),
				CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
				CorePairID:        uuid.FromStringOrNil(testCorePairUUID),
			},
			want: &ganymede.CPUCorePinned{
				UUID:              uuid.FromStringOrNil(testCPUPinnedCoreUUID),
				CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
				CorePairID:        uuid.FromStringOrNil(testCorePairUUID),
			},
			err: false,
		},
	}

	for _, test := range tests {
		err := testDatastore.PutPinnedCore(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		got, err := getPinnedCoreFromSQL(testDB, test.input.UUID)
		if err != nil {
			t.Fatalf("failed to get pinnedn core from SQL: %+v", err)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_DeletePinnedCore(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}
	if err := testDatastore.PutCPUPinningGroup(context.Background(), ganymede.CPUPinningGroup{
		UUID:         uuid.FromStringOrNil(testCPUPinningGroupUUID),
		Name:         testCPUPinningGroupName,
		HypervisorID: hypervisorID,
		CountCore:    4,
	}); err != nil {
		t.Fatalf("failed to put cpu pinning group: %+v", err)
	}
	nodes := []ganymede.NUMANode{
		ganymede.NUMANode{
			UUID: uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs: []ganymede.CorePair{
				ganymede.CorePair{
					UUID:         uuid.FromStringOrNil(testCorePairUUID),
					PhysicalCore: 0,
					LogicalCore:  41,
					NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
				},
			},
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 40,
			LogicalCoreMin:  41,
			LogicalCoreMax:  80,
			HypervisorID:    hypervisorID,
		},
	}
	if err = testDatastore.PutHypervisorNUMANode(context.Background(), nodes, hypervisorID); err != nil {
		t.Fatalf("failed to put hypervisor numa node: %+v", err)
	}

	pinned := ganymede.CPUCorePinned{
		UUID:              uuid.FromStringOrNil(testCPUPinnedCoreUUID),
		CPUPinningGroupID: uuid.FromStringOrNil(testCPUPinningGroupUUID),
		CorePairID:        uuid.FromStringOrNil(testCorePairUUID),
	}
	if err := testDatastore.PutPinnedCore(context.Background(), pinned); err != nil {
		t.Fatalf("failed to put pinned core: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.CPUCorePinned
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testCPUPinnedCoreUUID),
			want:  nil,
			err:   false,
		},
	}

	for _, test := range tests {
		err := testDatastore.DeletePinnedCore(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}

		got, err := getPinnedCoreFromSQL(testDB, uuid.FromStringOrNil(testCPUPinningGroupUUID))
		if test.want != got {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func getCPUPinningGroupFromSQL(testDB *sqlx.DB, cpgID uuid.UUID) (*ganymede.CPUPinningGroup, error) {
	var cpg ganymede.CPUPinningGroup
	query := `SELECT uuid, name, hypervisor_id, count_of_core, created_at, updated_at FROM cpu_pinning_group WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Get(&cpg, cpgID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &cpg, nil
}

func getPinnedCoreFromSQL(testDB *sqlx.DB, ccpID uuid.UUID) (*ganymede.CPUCorePinned, error) {
	var ccp ganymede.CPUCorePinned
	query := `SELECT uuid, pinning_group_id, hypervisor_cpu_pair_id, created_at, updated_at FROM cpu_core_pinned WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Get(&ccp, ccpID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &ccp, nil
}

func setFakeTimeNUMANodes(nodes []ganymede.NUMANode) []ganymede.NUMANode {
	var result []ganymede.NUMANode
	for _, node := range nodes {
		node.CreatedAt = time.Time{}
		node.UpdatedAt = time.Time{}

		var r []ganymede.CorePair
		for _, cp := range node.CorePairs {
			cp.CreatedAt = time.Time{}
			cp.UpdatedAt = time.Time{}

			r = append(r, cp)
		}

		node.CorePairs = r
		result = append(result, node)
	}

	return result
}
