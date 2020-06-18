package mysql

import (
	"context"
	"fmt"
	"net"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/ipam"
)

// CreateLease create a lease
func (m *MySQL) CreateLease(ctx context.Context, lease ipam.Lease) (*ipam.Lease, error) {
	query := `INSERT INTO lease(mac_address, address_id) VALUES (?, ?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, lease.MacAddress, lease.AddressID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &lease, nil
}

// GetLeaseByMACAddress retrieves lease according to the mac given
func (m *MySQL) GetLeaseByMACAddress(ctx context.Context, mac net.HardwareAddr) (*ipam.Lease, error) {
	query := `SELECT mac_address, address_id, created_at, updated_at FROM lease WHERE mac_address = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var lease ipam.Lease
	if err := stmt.GetContext(ctx, &lease, types.HardwareAddr(mac)); err != nil {
		return nil, fmt.Errorf("failed to get lease: %w", err)
	}

	return &lease, nil
}

// ListLease retrieves all leases
func (m *MySQL) ListLease(ctx context.Context) ([]ipam.Lease, error) {
	query := `SELECT mac_address, address_id, created_at, updated_at FROM lease`

	var leases []ipam.Lease
	if err := m.Conn.SelectContext(ctx, &leases, query); err != nil {
		return nil, fmt.Errorf("failed to get lease list: %w", err)
	}

	return leases, nil
}

// DeleteLease deletes a lease
func (m *MySQL) DeleteLease(ctx context.Context, mac net.HardwareAddr) error {
	query := `DELETE FROM lease WHERE mac_address = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, types.HardwareAddr(mac))
	if err != nil {
		return fmt.Errorf("failed to delete lease: %w", err)
	}

	return nil
}
