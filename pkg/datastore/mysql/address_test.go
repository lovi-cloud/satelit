package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/internal/testutils"
	"github.com/lovi-cloud/satelit/pkg/ipam"
)

func TestMySQL_CreateAddress(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input ipam.Address
		want  ipam.Address
		err   bool
	}{
		{
			input: ipam.Address{
				UUID:     uuid.FromStringOrNil(testAddressID),
				IP:       parseIP("192.168.2.100"),
				SubnetID: uuid.FromStringOrNil(testSubnetID),
			},
			want: ipam.Address{
				UUID:     uuid.FromStringOrNil(testAddressID),
				IP:       parseIP("192.168.2.100"),
				SubnetID: uuid.FromStringOrNil(testSubnetID),
			},
			err: false,
		},
	}
	for _, test := range tests {
		address, err := testDatastore.CreateAddress(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to create address: %+v", err)
		}
		got, err := getAddressFromSQL(testDB, address.UUID)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if ok, _ := testutils.CompareStruct(test.want, got); !ok {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestMySQL_GetAddressByID(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	_, err = testDatastore.CreateAddress(context.Background(), testAddress)
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  ipam.Address
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testAddressID),
			want:  testAddress,
			err:   false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetAddressByID(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if ok, _ := testutils.CompareStruct(test.want, got); !ok {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestMySQL_ListAddressBySubnetID(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	_, err = testDatastore.CreateAddress(context.Background(), testAddress)
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  []ipam.Address
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testSubnetID),
			want: []ipam.Address{
				testAddress,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.ListAddressBySubnetID(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if ok, _ := testutils.CompareStruct(test.want[0], got[0]); !ok {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestMySQL_DeleteAddress(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	_, err = testDatastore.CreateAddress(context.Background(), testAddress)
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ipam.Address
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testAddressID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DeleteAddress(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to delete address: %+v", err)
		}
		got, err := getAddressFromSQL(testDB, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if test.want != got {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func getAddressFromSQL(testDB *sqlx.DB, uuid uuid.UUID) (*ipam.Address, error) {
	var a ipam.Address
	query := `SELECT uuid, ip, subnet_id, created_at, updated_at FROM address WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare: %w", err)
	}
	err = stmt.Get(&a, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	return &a, nil
}
