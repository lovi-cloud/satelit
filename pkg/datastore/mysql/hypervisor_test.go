package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/go-test/deep"

	"github.com/jmoiron/sqlx"

	"github.com/lovi-cloud/satelit/internal/testutils"

	"github.com/lovi-cloud/satelit/pkg/ganymede"
)

const (
	testIQN      = "dummy-iqn"
	testHostname = "hv0001"
)

func TestMySQL_GetHypervisor(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	id, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}

	tests := []struct {
		input int
		want  *ganymede.HyperVisor
		err   bool
	}{
		{
			input: id,
			want: &ganymede.HyperVisor{
				ID:       id,
				IQN:      testIQN,
				Hostname: testHostname,
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetHypervisor(context.Background(), test.input)
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

func TestMySQL_GetHypervisorByHostname(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}

	tests := []struct {
		input string
		want  *ganymede.HyperVisor
		err   bool
	}{
		{
			input: testHostname,
			want: &ganymede.HyperVisor{
				ID:       1,
				IQN:      testIQN,
				Hostname: testHostname,
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetHypervisorByHostname(context.Background(), test.input)
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

func TestMySQL_PutHypervisor(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	tests := []struct {
		input []string // iqn, hostname
		want  *ganymede.HyperVisor
		err   bool
	}{
		{
			input: []string{testIQN, testHostname},
			want: &ganymede.HyperVisor{
				ID:       1,
				IQN:      testIQN,
				Hostname: testHostname,
			},
			err: false,
		},
	}

	for _, test := range tests {
		id, err := testDatastore.PutHypervisor(context.Background(), test.input[0], test.input[1])
		if err != nil {
			t.Fatalf("failed to put hypervisor: %+v", err)
		}
		got, err := getHypervisorFromSQL(testDB, id)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_PutHypervisorNUMANode(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	hypervisorID, err := testDatastore.PutHypervisor(context.Background(), testIQN, testHostname)
	if err != nil {
		t.Fatalf("failed to put hypervisor: %+v", err)
	}

	tests := []struct {
		input []ganymede.NUMANode
		want  []ganymede.NUMANode
		err   bool
	}{
		{
			input: []ganymede.NUMANode{
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
			},
			want: []ganymede.NUMANode{
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
			},
			err: false,
		},
	}

	for _, test := range tests {
		err := testDatastore.PutHypervisorNUMANode(context.Background(), test.input, hypervisorID)
		if err != nil {
			t.Fatalf("failed to put hypervisor numa node: %+v", err)
		}

		got, err := getNUMANodeFromSQL(testDB, test.input[0].UUID)
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

func getHypervisorFromSQL(testDB *sqlx.DB, hypervisorID int) (*ganymede.HyperVisor, error) {
	var hv ganymede.HyperVisor
	query := `SELECT id, iqn, hostname FROM hypervisor WHERE id = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Get(&hv, hypervisorID); err != nil {
		return nil, fmt.Errorf("failed to get hypervisor: %w", err)
	}

	return &hv, nil
}

func getNUMANodeFromSQL(testDB *sqlx.DB, numaNodeID uuid.UUID) ([]ganymede.NUMANode, error) {
	var nodes []ganymede.NUMANode
	query := `SELECT uuid, physical_core_min, physical_core_max, logical_core_min, logical_core_max, hypervisor_id FROM hypervisor_numa_node WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Select(&nodes, numaNodeID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute select numa node query: %w", err)
	}

	for i, node := range nodes {
		var pairs []ganymede.CorePair
		query := `SELECT uuid, numa_node_id, physical_core_number, logical_core_number FROM hypervisor_cpu_pair WHERE numa_node_id = ?`
		stmt, err := testDB.Preparex(query)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement: %w", err)
		}
		if err := stmt.Select(&pairs, node.UUID.String()); err != nil {
			return nil, fmt.Errorf("failed to execute select core pair query: %w", err)
		}

		for _, p := range pairs {
			nodes[i].CorePairs = append(nodes[i].CorePairs, p)
		}
	}

	return nodes, nil
}
