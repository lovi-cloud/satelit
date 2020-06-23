package mysql

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/ipam"
)

// CreateSubnet create a subnet
func (m *MySQL) CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*uuid.UUID, error) {
	query := `INSERT INTO subnet(uuid, name, network, start, end, gateway, dns_server, metadata_server) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	u := uuid.NewV4()
	_, err = stmt.ExecContext(ctx, u, subnet.Name, subnet.Network, subnet.Start, subnet.End, subnet.Gateway, subnet.DNSServer, subnet.MetadataServer)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &u, nil
}

// GetSubnetByID retrieves address according to the id given
func (m *MySQL) GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error) {
	query := `SELECT uuid, name, network, start, end, gateway, dns_server, metadata_server, created_at, updated_at FROM subnet WHERE uuid = UUID_TO_BIN(?)`
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
	query := `SELECT uuid, name, network, start, end, gateway, dns_server, metadata_server, created_at, updated_at FROM subnet`

	var subnets []ipam.Subnet
	if err := m.Conn.SelectContext(ctx, &subnets, query); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get subnet list: %w", err)
	} else if err == sql.ErrNoRows {
		return subnets, nil
	}

	return subnets, nil
}

// DeleteSubnet deletes a subnet
func (m *MySQL) DeleteSubnet(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM subnet WHERE uuid = UUID_TO_BIN(?)`
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
	query := `SELECT uuid, ip, subnet_id, created_at, updated_at FROM address WHERE uuid = UUID_TO_BIN(?)`
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
	query := `SELECT uuid, ip, subnet_id, created_at, updated_at FROM address WHERE subnet_id = UUID_TO_BIN(?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var addresses []ipam.Address
	if err := stmt.SelectContext(ctx, &addresses, subnetID); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get address list: %w", err)
	} else if err == sql.ErrNoRows {
		return addresses, nil
	}

	return addresses, nil
}

// DeleteAddress deletes address
func (m *MySQL) DeleteAddress(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM address WHERE uuid = UUID_TO_BIN(?)`
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
