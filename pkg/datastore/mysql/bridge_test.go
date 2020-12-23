package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/internal/testutils"
	"github.com/lovi-cloud/satelit/pkg/ganymede"
)

const (
	testBridgeID = "d09feb88-f30d-4b45-99cf-c4b1cdb22110"
)

func TestMySQL_CreateBridge(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	tests := []struct {
		input ganymede.Bridge
		want  *ganymede.Bridge
		err   bool
	}{
		{
			input: ganymede.Bridge{
				UUID:   uuid.FromStringOrNil(testBridgeID),
				VLANID: 1000,
				Name:   "testbr000",
			},
			want: &ganymede.Bridge{
				UUID:   uuid.FromStringOrNil(testBridgeID),
				VLANID: 1000,
				Name:   "testbr000",
			},
			err: false,
		},
	}
	for _, test := range tests {
		_, err := testDatastore.CreateBridge(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to create subnet: %+v", err)
		}
		got, err := getBridgeFromSQL(testDB, test.input.UUID)
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

func TestMySQL_GetBridge(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateBridge(context.Background(), ganymede.Bridge{
		UUID:   uuid.FromStringOrNil(testBridgeID),
		VLANID: 1000,
		Name:   "testbr1000",
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.Bridge
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testBridgeID),
			want: &ganymede.Bridge{
				UUID:   uuid.FromStringOrNil(testBridgeID),
				VLANID: 1000,
				Name:   "testbr1000",
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetBridge(context.Background(), test.input)
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

func TestMySQL_ListBridge(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateBridge(context.Background(), ganymede.Bridge{
		UUID:   uuid.FromStringOrNil(testBridgeID),
		VLANID: 1000,
		Name:   "testbr1000",
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		want []ganymede.Bridge
		err  bool
	}{
		{
			want: []ganymede.Bridge{
				{
					UUID:   uuid.FromStringOrNil(testBridgeID),
					VLANID: 1000,
					Name:   "testbr1000",
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.ListBridge(context.Background())
		if !test.err && err != nil {
			t.Fatalf("should not be error: %+v", err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error but not:")
		}
		if got != nil {
			for i := 0; i < len(got); i++ {
				got[i].CreatedAt = time.Time{}
				got[i].UpdatedAt = time.Time{}
			}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_DeleteBridge(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	_, err := testDatastore.CreateBridge(context.Background(), ganymede.Bridge{
		UUID:   uuid.FromStringOrNil(testBridgeID),
		VLANID: 1000,
		Name:   "testbr1000",
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.Bridge
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testBridgeID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DeleteBridge(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failedto delete bridge: %+v", err)
		}
		got, err := getBridgeFromSQL(testDB, test.input)
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

func getBridgeFromSQL(testDB *sqlx.DB, bridgeID uuid.UUID) (*ganymede.Bridge, error) {
	query := `SELECT uuid, vlan_id, name, created_at, updated_at FROM bridge where uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var b ganymede.Bridge
	err = stmt.Get(&b, bridgeID)
	if err != nil {
		return nil, err
	}

	return &b, nil
}
