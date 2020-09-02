package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/jmoiron/sqlx"

	"github.com/whywaita/satelit/internal/testutils"

	"github.com/whywaita/satelit/pkg/ganymede"
)

const (
	testIQN      = "dummy-iqn"
	testHostname = "hv0001"
)

func TestMySQL_GetHypervisorByHostname(t *testing.T) {
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
				ID:       1,
				IQN:      testIQN,
				Hostname: testHostname,
			},
		},
	}

	for _, test := range tests {
		got, err := testDatastore.GetHypervisorByHostname(context.Background(), testHostname)
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

func getHypervisorFromSQL(testDB *sqlx.DB, hypervisorID int) (*ganymede.HyperVisor, error) {
	var hv ganymede.HyperVisor
	query := `SELECT id, iqn, hostname FROM hypervisor WHERE id = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare: %w", err)
	}
	if err := stmt.Get(&hv, hypervisorID); err != nil {
		return nil, fmt.Errorf("failed to get hypervisor: %w", err)
	}

	return &hv, nil
}
