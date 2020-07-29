package mysql

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/ipam"
)

// CreateLease create a lease
func (m *MySQL) CreateLease(ctx context.Context, lease ipam.Lease) (*ipam.Lease, error) {
	query := `INSERT INTO lease(uuid, mac_address, address_id) VALUES (?, ?, ?)`
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

// GetLeaseByID retrieves lease according to the mac given
func (m *MySQL) GetLeaseByID(ctx context.Context, leaseID uuid.UUID) (*ipam.Lease, error) {
	query := `SELECT uuid, mac_address, address_id, created_at, updated_at FROM lease WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var lease ipam.Lease
	if err := stmt.GetContext(ctx, &lease, leaseID); err != nil {
		return nil, fmt.Errorf("failed to get lease: %w", err)
	}

	return &lease, nil
}

// GetDHCPLeaseByMACAddress retrieves DHCPLease according to the mac given
func (m *MySQL) GetDHCPLeaseByMACAddress(ctx context.Context, mac types.HardwareAddr) (*ipam.DHCPLease, error) {
	query := `select mac_address, ip, network, gateway, dns_server, metadata_server from lease lef join address on address_id = address.uuid left join subnet on subnet_id = subnet.uuid where mac_address = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	var dhcp ipam.DHCPLease
	if err := stmt.GetContext(ctx, &dhcp, mac); err != nil {
		return nil, fmt.Errorf("failed to get DHCP lease: %w", err)
	}

	return &dhcp, nil
}

// ListLease retrieves all leases
func (m *MySQL) ListLease(ctx context.Context) ([]ipam.Lease, error) {
	query := `SELECT uuid, mac_address, address_id, created_at, updated_at FROM lease`

	var leases []ipam.Lease
	if err := m.Conn.SelectContext(ctx, &leases, query); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get lease list: %w", err)
	} else if err == sql.ErrNoRows {
		return leases, nil
	}

	return leases, nil
}

// DeleteLease deletes a lease
func (m *MySQL) DeleteLease(ctx context.Context, leaseID int) error {
	query := `DELETE FROM lease WHERE uuid = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, leaseID)
	if err != nil {
		return fmt.Errorf("failed to delete lease: %w", err)
	}

	return nil
}
