package ipam

import (
	"bytes"
	"context"
	"net"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/datastore/memory"
	"github.com/whywaita/satelit/pkg/ipam"
)

func TestServerCreateSubnet(t *testing.T) {
	server := New(memory.New())

	tests := []struct {
		name           string
		vlan           uint32
		prefix         string
		start          string
		end            string
		gateway        string
		dnsServer      string
		metadataServer string
		want           *ipam.Subnet
		err            bool
	}{
		{
			name:           "test000",
			vlan:           1000,
			prefix:         "10.0.0.0/24",
			start:          "10.0.0.33",
			end:            "10.0.0.100",
			gateway:        "10.0.0.1",
			dnsServer:      "8.8.8.8",
			metadataServer: "10.0.0.15",
			want: &ipam.Subnet{
				Name:           "test000",
				VLANID:         1000,
				Network:        parseCIDR("10.0.0.0/24"),
				Start:          parseIP("10.0.0.33"),
				End:            parseIP("10.0.0.100"),
				Gateway:        parseIPPointer("10.0.0.1"),
				DNSServer:      parseIPPointer("8.8.8.8"),
				MetadataServer: parseIPPointer("10.0.0.15"),
			},
			err: false,
		},
		{
			name:   "test001",
			vlan:   1001,
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.100",
			want: &ipam.Subnet{
				Name:           "test001",
				VLANID:         1001,
				Network:        parseCIDR("192.168.0.0/25"),
				Start:          parseIP("192.168.0.33"),
				End:            parseIP("192.168.0.100"),
				Gateway:        nil,
				DNSServer:      nil,
				MetadataServer: nil,
			},
			err: false,
		},
		{
			name:   "test002",
			vlan:   1002,
			prefix: "192.168.0.0/25",
			start:  "192.168.1.33",
			end:    "192.168.1.100",
			want:   nil,
			err:    true,
		},
		{
			name:   "test003",
			vlan:   1003,
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.33",
			want:   nil,
			err:    true,
		},
		{
			name:   "test004",
			vlan:   1004,
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.32",
			want:   nil,
			err:    true,
		},
		{
			name:   "test005",
			vlan:   1006,
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.254",
			want:   nil,
			err:    true,
		},
		{
			name:    "test006",
			vlan:    1006,
			prefix:  "192.168.0.0/25",
			start:   "192.168.0.33",
			end:     "192.168.0.34",
			gateway: "10.0.0.1",
			want:    nil,
			err:     true,
		},
		{
			name:      "test007",
			vlan:      1007,
			prefix:    "192.168.0.0/25",
			start:     "192.168.0.33",
			end:       "192.168.0.34",
			dnsServer: "example.com",
			want:      nil,
			err:       true,
		},
		{
			name:           "test008",
			vlan:           1008,
			prefix:         "192.168.0.0/25",
			start:          "192.168.0.33",
			end:            "192.168.0.34",
			metadataServer: "1.1.1.1",
			want:           nil,
			err:            true,
		},
	}
	for _, test := range tests {
		got, err := server.CreateSubnet(context.Background(), test.name, test.vlan, test.prefix, test.start, test.end, test.gateway, test.dnsServer, test.metadataServer)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %v but: %v", test.name, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.name)
		}
		if !subnetEqual(got, test.want) {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestServerCreateAddress(t *testing.T) {
	server := New(memory.New())

	subnet, err := server.CreateSubnet(context.Background(), "test000", 1000, "10.0.0.0/24", "10.0.0.100", "10.0.0.102", "10.0.0.1", "8.8.8.8", "10.0.0.15")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		subnetID uuid.UUID
		want     *ipam.Address
		err      bool
	}{
		{
			subnetID: subnet.UUID,
			want: &ipam.Address{
				IP:       parseIP("10.0.0.100"),
				SubnetID: subnet.UUID,
			},
			err: false,
		},
		{
			subnetID: subnet.UUID,
			want: &ipam.Address{
				IP:       parseIP("10.0.0.101"),
				SubnetID: subnet.UUID,
			},
			err: false,
		},
		{
			subnetID: subnet.UUID,
			want: &ipam.Address{
				IP:       parseIP("10.0.0.102"),
				SubnetID: subnet.UUID,
			},
			err: false,
		},
		{
			subnetID: subnet.UUID,
			want:     nil,
			err:      true,
		},
		{
			subnetID: uuid.UUID{},
			want:     nil,
			err:      true,
		},
	}
	for _, test := range tests {
		got, err := server.CreateAddress(context.Background(), test.subnetID)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %v but: %v", test.subnetID, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.subnetID)
		}
		if !addressEqual(got, test.want) {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}

}

func parseIP(s string) types.IP {
	return types.IP(net.ParseIP(s))
}

func parseIPPointer(s string) *types.IP {
	i := parseIP(s)
	return &i
}

func parseCIDR(s string) types.IPNet {
	_, net, _ := net.ParseCIDR(s)
	return types.IPNet(*net)
}

func subnetEqual(a, b *ipam.Subnet) bool {
	if a == b {
		return true
	}
	if a.Name != b.Name {
		return false
	}
	if bytes.Compare(a.Network.IP, b.Network.IP) != 0 {
		return false
	}
	if bytes.Compare(a.Network.Mask, b.Network.Mask) != 0 {
		return false
	}
	if bytes.Compare(a.Start, b.Start) != 0 {
		return false
	}
	if bytes.Compare(a.End, b.End) != 0 {
		return false
	}
	return true
}

func addressEqual(a, b *ipam.Address) bool {
	if a == b {
		return true
	}
	if bytes.Compare(a.IP, b.IP) != 0 {
		return false
	}
	if a.SubnetID != b.SubnetID {
		return false
	}
	return true
}
