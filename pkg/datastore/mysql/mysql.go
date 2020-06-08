package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/whywaita/satelit/api"

	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A MySQL is backend of datastore by MySQL Server
type MySQL struct {
	Conn *sqlx.DB
}

// New create MySQL datastore
func New(c *config.MySQLConfig) (*MySQL, error) {
	dsn := c.DSN + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect : %w", err)
	}

	conn.SetMaxIdleConns(c.MaxIdleConn)
	conn.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetimeSecond) * time.Second)

	return &MySQL{
		Conn: conn,
	}, nil
}

// GetIQN get IQN from MySQL
func (m *MySQL) GetIQN(ctx context.Context, hostname string) (string, error) {
	var iqn string

	query := fmt.Sprintf(`SELECT iqn FROM hypervisor WHERE hostname = "%s"`, hostname)
	err := m.Conn.Get(&iqn, query)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		break
	case err != nil:
		return "", fmt.Errorf("failed to exec query (hostname: %v): %w", hostname, err)
	}

	if iqn != "" {
		return iqn, nil
	}

	// not found in mysql
	resp, err := teleskop.GetClient(hostname).GetISCSIQualifiedName(ctx, &pb.GetISCSIQualifiedNameRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get IQN from Teleskop (host: %v): %w", hostname, err)
	}
	return resp.Iqn, nil
}

// GetImage return image object
func (m *MySQL) GetImage(imageID string) (*europa.BaseImage, error) {
	var image europa.BaseImage

	query := fmt.Sprintf(`SELECT * FROM image WHERE uuid = "%s"`, imageID)
	err := m.Conn.Get(&image, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &image, nil
}

// GetImages return all images
func (m *MySQL) GetImages() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	query := fmt.Sprintf("SELECT * FROM image")
	err := m.Conn.Select(&images, query)
	if err != nil {
		return nil, fmt.Errorf("failed to SELCT image table: %w", err)
	}

	return images, nil
}

// PutImage write image record
func (m *MySQL) PutImage(image europa.BaseImage) error {
	query := `INSERT INTO image(uuid, name, volume_id, description) VALUES (UUID_TO_BIN(?), ?, ?, ?)`
	_, err := m.Conn.Exec(query, image.UUID, image.Name, image.CacheVolumeID, image.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// DeleteImage delete image record
func (m *MySQL) DeleteImage(imageID string) error {
	query := fmt.Sprintf(`DELETE FROM image WHERE uuid = "%s"`, imageID)
	_, err := m.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}

// CreateSubnet create a subnet
func (m *MySQL) CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*uuid.UUID, error) {
	query := `INSERT INTO subnet(uuid, name, network, start, end) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	u := uuid.NewV4()
	_, err = stmt.ExecContext(ctx, u, subnet.Name, subnet.Network, subnet.Start, subnet.End)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &u, nil
}

// GetSubnetByID retrieves address according to the id given
func (m *MySQL) GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error) {
	query := `SELECT uuid, name, network, start, end, created_at, updated_at FROM subnet WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var subnet ipam.Subnet
	if err := stmt.GetContext(ctx, &subnet, uuid); err != nil {
		return nil, fmt.Errorf("failed to get subnet: %w", err)
	}

	return &subnet, nil
}

// ListSubnet retrieves all subnets
func (m *MySQL) ListSubnet(ctx context.Context) ([]ipam.Subnet, error) {
	query := `SELECT uuid, name, network, start, end, created_at, updated_at FROM subnet`

	var subnets []ipam.Subnet
	if err := m.Conn.SelectContext(ctx, &subnets, query); err != nil {
		return nil, fmt.Errorf("failed to get subnet list: %w", err)
	}

	return subnets, nil
}

// DeleteSubnet deletes a subnet
func (m *MySQL) DeleteSubnet(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM subnet WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete subnet: %w", err)
	}

	return nil
}

// CreateAddress create a address
func (m *MySQL) CreateAddress(ctx context.Context, address ipam.Address) (*uuid.UUID, error) {
	query := `INSERT INTO address(uuid, ip, subnet_id) VALUES (UUID_TO_BIN(?), ?, UUID_TO_BIN(?))`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	u := uuid.NewV4()
	_, err = stmt.ExecContext(ctx, u, address.IP, address.SubnetID)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return &u, nil
}

// GetAddressByID retrieves address according to the id given
func (m *MySQL) GetAddressByID(ctx context.Context, uuid uuid.UUID) (*ipam.Address, error) {
	query := `SELECT uuid, ip, subnet_id, created_at, updated_at FROM address WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var address ipam.Address
	if err := stmt.GetContext(ctx, &address, uuid); err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return &address, nil
}

// ListAddressBySubnetID retrieves all address according to the subnetID given.
func (m *MySQL) ListAddressBySubnetID(ctx context.Context, subnetID uuid.UUID) ([]ipam.Address, error) {
	query := `SELECT uuid, ip, subnet_id, created_at, updated_at FROM address WHERE subnet_id = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var addresses []ipam.Address
	if err := stmt.SelectContext(ctx, &addresses, subnetID); err != nil {
		return nil, fmt.Errorf("failed to get address list: %w", err)
	}

	return addresses, nil
}

// DeleteAddress deletes address
func (m *MySQL) DeleteAddress(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM address WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

// GetVirtualMachine return virtual machine record
func (m *MySQL) GetVirtualMachine(vmUUID string) (*ganymede.VirtualMachine, error) {
	var vm ganymede.VirtualMachine
	query := fmt.Sprintf(`SELECT 
BIN_TO_UUID(uuid),
name,
vcpus,
memory_kib,
hypervisor_name WHERE uuid = %s`, vmUUID)
	err := m.Conn.Get(&vm, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &vm, nil
}

// PutVirtualMachine write virtual machine record
func (m *MySQL) PutVirtualMachine(vm ganymede.VirtualMachine) error {
	query := `INSERT INTO virtual_machine(name, uuid, vcpus, memory_kib, hypervisor_name) VALUES (?, UUID_TO_BIN(?), ?, ?, ?)`
	_, err := m.Conn.Exec(query, vm.Name, vm.UUID, vm.Vcpus, vm.MemoryKiB, vm.HypervisorName)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
