package mysql

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/satelit/pkg/ipam"
)

// CreateAddress create a address
func (m *MySQL) CreateAddress(ctx context.Context, address ipam.Address) (*ipam.Address, error) {
	query := `INSERT INTO address(uuid, ip, subnet_id) VALUES (?, ?, ?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, address.UUID, address.IP, address.SubnetID)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return &address, nil
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
	if err := stmt.SelectContext(ctx, &addresses, subnetID); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get address list: %w", err)
	} else if err == sql.ErrNoRows {
		return addresses, nil
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
