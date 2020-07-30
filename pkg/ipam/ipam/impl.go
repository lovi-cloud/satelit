package ipam

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/ipam"
)

type server struct {
	mutex     *sync.Mutex
	datastore datastore.Datastore
}

// New is return new IPAM interface
func New(d datastore.Datastore) ipam.IPAM {
	return &server{
		mutex:     &sync.Mutex{},
		datastore: d,
	}
}

// CreateSubnet create a subnet
func (s server) CreateSubnet(ctx context.Context, name string, vlanID uint32, prefix, start, end, gateway, dnsServer, metadataServer string) (*ipam.Subnet, error) {
	_, prefixNet, err := net.ParseCIDR(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prefix: %w", err)
	}

	startAddr := net.ParseIP(start)
	if startAddr == nil {
		return nil, fmt.Errorf("failed to paser start address")
	}
	if !prefixNet.Contains(startAddr) {
		return nil, fmt.Errorf("invalid start address")
	}

	endAddr := net.ParseIP(end)
	if endAddr == nil {
		return nil, fmt.Errorf("failed to parse end address")
	}
	if !prefixNet.Contains(endAddr) {
		return nil, fmt.Errorf("invalid end address")
	}

	if bytes.Compare(startAddr, endAddr) != -1 {
		return nil, fmt.Errorf("start address must be before the end address")
	}

	subnet := ipam.Subnet{
		UUID:           uuid.NewV4(),
		Name:           name,
		VLANID:         vlanID,
		Network:        types.IPNet(*prefixNet),
		Start:          types.IP(startAddr),
		End:            types.IP(endAddr),
		Gateway:        nil,
		DNSServer:      nil,
		MetadataServer: nil,
	}

	if gateway != "" {
		gwAddr := net.ParseIP(gateway)
		if gwAddr == nil {
			return nil, fmt.Errorf("failed to parse gateway address")
		}
		if !prefixNet.Contains(gwAddr) {
			return nil, fmt.Errorf("invalid gateway address")
		}
		gw := types.IP(gwAddr)
		subnet.Gateway = &gw
	}

	if dnsServer != "" {
		dnsAddr := net.ParseIP(dnsServer)
		if dnsAddr == nil {
			return nil, fmt.Errorf("failed to parse DNS server address")
		}
		dns := types.IP(dnsAddr)
		subnet.DNSServer = &dns
	}

	if metadataServer != "" {
		mdAddr := net.ParseIP(metadataServer)
		if mdAddr == nil {
			return nil, fmt.Errorf("failed to parse metadata server address")
		}
		if !prefixNet.Contains(mdAddr) {
			return nil, fmt.Errorf("invalid metadata server address")
		}
		md := types.IP(mdAddr)
		subnet.MetadataServer = &md
	}

	return s.datastore.CreateSubnet(ctx, subnet)
}

// GetSubnet retrieves address according to the id given
func (s server) GetSubnet(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error) {
	return s.datastore.GetSubnetByID(ctx, uuid)
}

// ListSubnet retrieves all subnets
func (s server) ListSubnet(ctx context.Context) ([]ipam.Subnet, error) {
	return s.datastore.ListSubnet(ctx)
}

// DeleteSubnet deletes a subnet
func (s server) DeleteSubnet(ctx context.Context, uuid uuid.UUID) error {
	return s.datastore.DeleteSubnet(ctx, uuid)
}

// CreateAddress create a address
func (s server) CreateAddress(ctx context.Context, subnetID uuid.UUID) (*ipam.Address, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subnet, err := s.datastore.GetSubnetByID(ctx, subnetID)
	if err != nil {
		return nil, err
	}

	addresses, err := s.datastore.ListAddressBySubnetID(ctx, subnetID)
	if err != nil {
		return nil, err
	}

	start := net.IP(subnet.Start)
	end := net.IP(subnet.End)
	var addr *net.IP
	network := net.IPNet(subnet.Network)
	for network.Contains(start) && bytes.Compare(start, end) < 1 {
		isExist := false
		for _, a := range addresses {
			if bytes.Compare(a.IP, start) == 0 {
				isExist = true
				break
			}
		}
		if !isExist {
			addr = &start
			break
		}
		start = getNextAddress(start)
	}
	if addr == nil {
		return nil, fmt.Errorf("failed to get address")
	}

	address := ipam.Address{
		UUID:     uuid.NewV4(),
		IP:       types.IP(*addr),
		SubnetID: subnetID,
	}
	return s.datastore.CreateAddress(ctx, address)
}

// GetAddress retrieves address according to the id given
func (s server) GetAddress(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error) {
	return s.datastore.GetAddressByID(ctx, uuid)
}

// ListAddressBySubnetID retrieves all address according to the subnetID given.
func (s server) ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error) {
	return s.datastore.ListAddressBySubnetID(ctx, subnetID)
}

// DeleteAddress deletes address
func (s server) DeleteAddress(ctx context.Context, uuid uuid.UUID) error {
	return s.datastore.DeleteAddress(ctx, uuid)
}

func (s server) CreateLease(ctx context.Context, addressID uuid.UUID) (*ipam.Lease, error) {
	mac, err := generateNewMACAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new MAC Address: %w", err)
	}

	lease := ipam.Lease{
		UUID:       uuid.NewV4(),
		MacAddress: types.HardwareAddr(*mac),
		AddressID:  addressID,
	}
	return s.datastore.CreateLease(ctx, lease)
}

func (s server) GetLease(ctx context.Context, leaseID uuid.UUID) (*ipam.Lease, error) {
	return s.datastore.GetLeaseByID(ctx, leaseID)
}

func (s server) ListLease(ctx context.Context) ([]ipam.Lease, error) {
	return s.datastore.ListLease(ctx)
}

func (s server) DeleteLease(ctx context.Context, leaseID uuid.UUID) error {
	return s.datastore.DeleteLease(ctx, leaseID)
}

func getNextAddress(ip net.IP) net.IP {
	a := net.ParseIP(ip.String())
	for i := len(a) - 1; i >= 0; i-- {
		a[i]++
		if a[i] > 0 {
			break
		}
	}
	return a
}

func generateNewMACAddress() (*net.HardwareAddr, error) {
	var mac net.HardwareAddr

	buff := make([]byte, 3)
	_, err := rand.Read(buff)
	if err != nil {
		return nil, err
	}
	mac = append(mac, byte(0xca), byte(0x03), byte(0x18), buff[0], buff[1], buff[2])

	return &mac, nil
}
