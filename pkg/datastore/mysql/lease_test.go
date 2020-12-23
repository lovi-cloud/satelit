package mysql_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/internal/mysql/types"
	"github.com/lovi-cloud/satelit/internal/testutils"
	"github.com/lovi-cloud/satelit/pkg/ipam"
)

const (
	testSubnetID  = "bba39c58-4af7-46aa-ab3d-eac7ab7b581b"
	testAddressID = "1d2f0d83-1508-4af0-bdcb-c52c34733923"
	testLeaseID   = "3e3e4f46-cd5d-4cf7-86a3-6b7af5d754e9"
	testMACAddr   = "ca:03:18:00:00:00"
)

var testSubnet = ipam.Subnet{
	UUID:           uuid.FromStringOrNil(testSubnetID),
	Name:           "lease-test",
	Network:        parseIPNet("192.168.2.0/24"),
	Start:          parseIP("192.168.2.100"),
	End:            parseIP("192.168.2.200"),
	Gateway:        parseIPp("192.168.2.1"),
	DNSServer:      parseIPp("192.168.2.2"),
	MetadataServer: parseIPp("192.168.2.3"),
}

var testAddress = ipam.Address{
	UUID:     uuid.FromStringOrNil(testAddressID),
	IP:       parseIP("192.168.2.100"),
	SubnetID: testSubnet.UUID,
}

var testLease = ipam.Lease{
	UUID:       uuid.FromStringOrNil(testLeaseID),
	MacAddress: parseMAC(testMACAddr),
	AddressID:  testAddress.UUID,
}

func TestMySQL_CreateLease(t *testing.T) {
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
		input ipam.Lease
		err   bool
	}{
		{
			input: testLease,
			err:   false,
		},
	}
	for _, test := range tests {
		_, err := testDatastore.CreateLease(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to create lease: %+v", err)
		}
		want, err := getLeaseFromSQL(testDB, test.input.UUID)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if ok, _ := testutils.CompareStruct(test.input, want); !ok {
			t.Fatalf("want %q, but %q:", test.input, want)
		}
	}
}

func TestMySQL_GetLeaseByMACAddress(t *testing.T) {
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
	_, err = testDatastore.CreateLease(context.Background(), testLease)
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input uuid.UUID
		want  ipam.Lease
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testLeaseID),
			want:  testLease,
			err:   false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetLeaseByID(context.Background(), test.input)
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

func TestMySQL_GetDHCPLeaseByMACAddress(t *testing.T) {
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
	_, err = testDatastore.CreateLease(context.Background(), testLease)
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input types.HardwareAddr
		want  ipam.DHCPLease
		err   bool
	}{
		{
			input: parseMAC(testMACAddr),
			want: ipam.DHCPLease{
				MacAddress:     testLease.MacAddress,
				IP:             testAddress.IP,
				Network:        testSubnet.Network,
				Gateway:        testSubnet.Gateway,
				DNSServer:      testSubnet.DNSServer,
				MetadataServer: testSubnet.MetadataServer,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetDHCPLeaseByMACAddress(context.Background(), test.input)
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

func TestMySQL_ListLease(t *testing.T) {
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
	_, err = testDatastore.CreateLease(context.Background(), testLease)
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input ipam.Lease
		want  []ipam.Lease
		err   bool
	}{
		{
			input: testLease,
			want: []ipam.Lease{
				testLease,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.ListLease(context.Background())
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

func TestMySQL_DeleteLease(t *testing.T) {
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
		want  *ipam.Lease
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testLeaseID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DeleteLease(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to delete lease: %+v", err)
		}
		got, err := getLeaseFromSQL(testDB, test.input)
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

func parseIP(s string) types.IP {
	return types.IP(net.ParseIP(s))
}

func parseIPp(s string) *types.IP {
	i := types.IP(net.ParseIP(s))
	return &i
}

func parseIPNet(s string) types.IPNet {
	_, ipNet, _ := net.ParseCIDR(s)
	return types.IPNet(*ipNet)
}

func parseMAC(s string) types.HardwareAddr {
	mac, _ := net.ParseMAC(s)
	return types.HardwareAddr(mac)
}

func getLeaseFromSQL(testDB *sqlx.DB, leaseID uuid.UUID) (*ipam.Lease, error) {
	var l ipam.Lease
	query := `SELECT uuid, mac_address, address_id, created_at, updated_at FROM lease WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare: %w", err)
	}
	err = stmt.Get(&l, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lease: %w", err)
	}
	return &l, nil
}
