package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/testutils"
	"github.com/whywaita/satelit/pkg/ipam"
)

func TestMySQL_CreateSubnet(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	tests := []struct {
		input ipam.Subnet
		want  ipam.Subnet
		err   bool
	}{
		{
			input: ipam.Subnet{
				UUID:           uuid.FromStringOrNil(testSubnetID),
				Name:           "test-subnet",
				Network:        parseIPNet("192.168.2.0/24"),
				Start:          parseIP("192.168.2.100"),
				End:            parseIP("192.168.2.200"),
				Gateway:        parseIPp("192.168.2.1"),
				DNSServer:      parseIPp("192.168.2.2"),
				MetadataServer: parseIPp("192.168.2.3"),
			},
			want: ipam.Subnet{
				UUID:           uuid.FromStringOrNil(testSubnetID),
				Name:           "test-subnet",
				Network:        parseIPNet("192.168.2.0/24"),
				Start:          parseIP("192.168.2.100"),
				End:            parseIP("192.168.2.200"),
				Gateway:        parseIPp("192.168.2.1"),
				DNSServer:      parseIPp("192.168.2.2"),
				MetadataServer: parseIPp("192.168.2.3"),
			},
			err: false,
		},
	}
	for _, test := range tests {
		subnet, err := testDatastore.CreateSubnet(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to create subnet: %+v", err)
		}
		got, err := getSubnetFromSQL(testDB, subnet.UUID)
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

func TestMySQL_GetSubnetByID(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  ipam.Subnet
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testSubnetID),
			want:  testSubnet,
			err:   false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetSubnetByID(context.Background(), test.input)
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

func TestMySQL_ListSubnet(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input ipam.Subnet
		want  []ipam.Subnet
		err   bool
	}{
		{
			input: testSubnet,
			want: []ipam.Subnet{
				testSubnet,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.ListSubnet(context.Background())
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

func TestMySQL_DeleteSubnet(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	_, err := testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ipam.Subnet
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testSubnetID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DeleteSubnet(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to create subnet: %+v", err)
		}
		got, err := getSubnetFromSQL(testDB, test.input)
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

func getSubnetFromSQL(testDB *sqlx.DB, uuid uuid.UUID) (*ipam.Subnet, error) {
	var s ipam.Subnet
	query := `SELECT uuid, name, network, start, end, gateway, dns_server, metadata_server, created_at, updated_at FROM subnet WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare: %w", err)
	}
	err = stmt.Get(&s, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet: %w", err)
	}
	return &s, nil
}
