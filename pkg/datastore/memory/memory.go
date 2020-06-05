package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A Memory is on memory datastore for testing.
type Memory struct {
	mutex *sync.Mutex

	images    map[string]europa.BaseImage
	subnets   map[uuid.UUID]ipam.Subnet
	addresses map[uuid.UUID]ipam.Address
}

// New create Memory
func New() *Memory {
	return &Memory{
		mutex:     &sync.Mutex{},
		subnets:   map[uuid.UUID]ipam.Subnet{},
		addresses: map[uuid.UUID]ipam.Address{},
	}
}

// GetIQN return IQN from on memory
func (m *Memory) GetIQN(ctx context.Context, hostname string) (string, error) {
	return "dummy iqn", nil
}

// GetImage return image by id from on memory
func (m *Memory) GetImage(imageID string) (*europa.BaseImage, error) {
	m.mutex.Lock()
	i, ok := m.images[imageID]
	m.mutex.Unlock()

	if ok == false {
		return nil, errors.New("not found")
	}

	return &i, nil
}

// GetImages return all images
func (m *Memory) GetImages() ([]europa.BaseImage, error) {
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
	m.images[string(image.ID)] = image
	m.mutex.Unlock()
	return nil
}

// DeleteImage delete image
func (m *Memory) DeleteImage(imageID string) error {
	m.mutex.Lock()
	delete(m.images, imageID)
	m.mutex.Unlock()

	return nil
}

// CreateSubnet create a subnet
func (m *Memory) CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*uuid.UUID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	subnet.UUID = uuid.NewV4()
	m.subnets[subnet.UUID] = subnet

	return &subnet.UUID, nil
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
func (m *Memory) CreateAddress(ctx context.Context, address ipam.Address) (*uuid.UUID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	address.UUID = uuid.NewV4()
	m.addresses[address.UUID] = address

	return &address.UUID, nil
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
