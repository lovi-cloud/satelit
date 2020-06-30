package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A Memory is on memory datastore for testing.
type Memory struct {
	mutex *sync.Mutex

	volumes         map[string]europa.Volume
	images          map[uuid.UUID]europa.BaseImage
	subnets         map[uuid.UUID]ipam.Subnet
	addresses       map[uuid.UUID]ipam.Address
	leases          map[string]ipam.Lease
	virtualMachines map[uuid.UUID]ganymede.VirtualMachine
}

// New create Memory
func New() *Memory {
	return &Memory{
		mutex:           &sync.Mutex{},
		volumes:         map[string]europa.Volume{},
		images:          map[uuid.UUID]europa.BaseImage{},
		subnets:         map[uuid.UUID]ipam.Subnet{},
		addresses:       map[uuid.UUID]ipam.Address{},
		leases:          map[string]ipam.Lease{},
		virtualMachines: map[uuid.UUID]ganymede.VirtualMachine{},
	}
}

// GetIQN return IQN from on memory
func (m *Memory) GetIQN(ctx context.Context, hostname string) (string, error) {
	return "dummy iqn", nil
}

// GetImage return image by id from on memory
func (m *Memory) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	m.mutex.Lock()
	i, ok := m.images[imageID]
	m.mutex.Unlock()

	if ok == false {
		return nil, errors.New("not found")
	}

	return &i, nil
}

// ListImage retrieves all images
func (m *Memory) ListImage() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	m.mutex.Lock()
	for _, v := range m.images {
		images = append(images, v)
	}
	m.mutex.Unlock()

	return images, nil
}

// PutImage write image
func (m *Memory) PutImage(image europa.BaseImage) error {
	m.mutex.Lock()
	m.images[image.UUID] = image
	m.mutex.Unlock()
	return nil
}

// DeleteImage delete image
func (m *Memory) DeleteImage(imageID uuid.UUID) error {
	m.mutex.Lock()
	delete(m.images, imageID)
	m.mutex.Unlock()

	return nil
}

// ListVolume rerieves volumes
func (m *Memory) ListVolume(volumeIDs []string) ([]europa.Volume, error) {
	var vs []europa.Volume

	m.mutex.Lock()
	for _, volumeID := range volumeIDs {
		v, ok := m.volumes[volumeID]
		if ok {
			vs = append(vs, v)
		}
	}
	m.mutex.Unlock()

	return vs, nil
}

// GetVolume return volume
func (m *Memory) GetVolume(volumeID string) (*europa.Volume, error) {
	m.mutex.Lock()
	v, ok := m.volumes[volumeID]
	m.mutex.Unlock()

	if !ok {
		return nil, errors.New("not found")
	}

	return &v, nil
}

// PutVolume write volume
func (m *Memory) PutVolume(volume europa.Volume) error {
	m.mutex.Lock()
	m.volumes[volume.ID] = volume
	m.mutex.Unlock()

	return nil
}

// DeleteVolume delete volume
func (m *Memory) DeleteVolume(volumeID string) error {
	m.mutex.Lock()
	delete(m.volumes, volumeID)
	m.mutex.Unlock()

	return nil
}

// CreateSubnet create a subnet
func (m *Memory) CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*ipam.Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.subnets[subnet.UUID] = subnet

	return &subnet, nil
}

// GetSubnetByID retrieves address according to the id given
func (m *Memory) GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.subnets[uuid]
	if !ok {
		return nil, fmt.Errorf("failed to find subnet uuid=%s", uuid)
	}

	return &val, nil
}

// ListSubnet retrieves all subnets
func (m *Memory) ListSubnet(ctx context.Context) ([]ipam.Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var subnets []ipam.Subnet
	for _, subnet := range m.subnets {
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

// DeleteSubnet deletes a subnet
func (m *Memory) DeleteSubnet(ctx context.Context, uuid uuid.UUID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.subnets[uuid]
	if !ok {
		return fmt.Errorf("failed to find subnet uuid=%s", uuid)
	}
	delete(m.subnets, uuid)

	return nil
}

// CreateAddress create a address
func (m *Memory) CreateAddress(ctx context.Context, address ipam.Address) (*ipam.Address, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.addresses[address.UUID] = address

	return &address, nil
}

// GetAddressByID retrieves address according to the id given
func (m *Memory) GetAddressByID(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.addresses[uuid]
	if !ok {
		return nil, fmt.Errorf("failed to find address uuid=%s", uuid)
	}

	return &val, nil
}

// ListAddressBySubnetID retrieves all address according to the subnetID given.
func (m *Memory) ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var addresses []ipam.Address
	for _, address := range m.addresses {
		if address.SubnetID == subnetID {
			addresses = append(addresses, address)
		}
	}

	return addresses, nil
}

// DeleteAddress deletes address
func (m *Memory) DeleteAddress(ctx context.Context, uuid uuid.UUID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.addresses[uuid]
	if !ok {
		return fmt.Errorf("failed to find address uuid=%s", uuid)
	}
	delete(m.addresses, uuid)

	return nil
}

// CreateLease create a lease
func (m *Memory) CreateLease(ctx context.Context, lease ipam.Lease) (*ipam.Lease, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	lease.CreatedAt = now
	lease.UpdatedAt = now
	m.leases[lease.MacAddress.String()] = lease

	return &lease, nil
}

// GetLeaseByMACAddress retrieves lease according to the mac given
func (m *Memory) GetLeaseByMACAddress(ctx context.Context, mac types.HardwareAddr) (*ipam.Lease, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.leases[mac.String()]
	if !ok {
		return nil, fmt.Errorf("failed to find lease mac_address=%s", mac.String())
	}

	return &val, nil
}

// GetDHCPLeaseByMACAddress retrieves DHCPLease according to the mac given
func (m *Memory) GetDHCPLeaseByMACAddress(ctx context.Context, mac types.HardwareAddr) (*ipam.DHCPLease, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	lease, ok := m.leases[mac.String()]
	if !ok {
		return nil, fmt.Errorf("failed to find lease mac_address=%s", mac.String())
	}
	address, ok := m.addresses[lease.AddressID]
	if !ok {
		return nil, fmt.Errorf("failed to find address uuid=%s", lease.AddressID)
	}
	subnet, ok := m.subnets[address.SubnetID]
	if !ok {
		return nil, fmt.Errorf("failed to find subnet uuid=%s", address.SubnetID)
	}

	return &ipam.DHCPLease{
		MacAddress:     lease.MacAddress,
		IP:             address.IP,
		Gateway:        subnet.Gateway,
		DNSServer:      subnet.DNSServer,
		MetadataServer: subnet.MetadataServer,
	}, nil
}

// ListLease retrieves all leases
func (m *Memory) ListLease(ctx context.Context) ([]ipam.Lease, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var leases []ipam.Lease
	for _, lease := range m.leases {
		leases = append(leases, lease)
	}

	return leases, nil
}

// DeleteLease deletes a lease
func (m *Memory) DeleteLease(ctx context.Context, mac types.HardwareAddr) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.leases[mac.String()]
	if !ok {
		return fmt.Errorf("failed to find lease mac_address=%s", mac.String())
	}
	delete(m.leases, mac.String())

	return nil
}

// GetVirtualMachine return virtual machine record
func (m *Memory) GetVirtualMachine(vmUUID string) (*ganymede.VirtualMachine, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id, err := uuid.FromString(vmUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	val, ok := m.virtualMachines[id]
	if !ok {
		return nil, fmt.Errorf("failed to find virtual machine uuid=%s", id)
	}

	return &val, nil
}

// PutVirtualMachine write virtual machine record
func (m *Memory) PutVirtualMachine(vm ganymede.VirtualMachine) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.virtualMachines[vm.UUID] = vm

	return nil
}
