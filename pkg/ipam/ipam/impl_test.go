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
		name   string
		prefix string
		start  string
		end    string
		want   *ipam.Subnet
		err    bool
	}{
		{
			name:   "test000",
			prefix: "10.0.0.0/24",
			start:  "10.0.0.33",
			end:    "10.0.0.100",
			want: &ipam.Subnet{
				Name:    "test000",
				Network: parseCIDR("10.0.0.0/24"),
				Start:   parseIP("10.0.0.33"),
				End:     parseIP("10.0.0.100"),
			},
			err: false,
		},
		{
			name:   "test001",
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.100",
			want: &ipam.Subnet{
				Name:    "test001",
				Network: parseCIDR("192.168.0.0/25"),
				Start:   parseIP("192.168.0.33"),
				End:     parseIP("192.168.0.100"),
			},
			err: false,
		},
		{
			name:   "test002",
			prefix: "192.168.0.0/25",
			start:  "192.168.1.33",
			end:    "192.168.1.100",
			want:   nil,
			err:    true,
		},
		{
			name:   "test003",
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.33",
			want:   nil,
			err:    true,
		},
		{
			name:   "test004",
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.32",
			want:   nil,
			err:    true,
		},
		{
			name:   "test005",
			prefix: "192.168.0.0/25",
			start:  "192.168.0.33",
			end:    "192.168.0.254",
			want:   nil,
			err:    true,
		},
	}
	for _, test := range tests {
		got, err := server.CreateSubnet(context.Background(), test.name, test.prefix, test.start, test.end)
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

	subnet, err := server.CreateSubnet(context.Background(), "test000", "10.0.0.0/24", "10.0.0.100", "10.0.0.102")
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
