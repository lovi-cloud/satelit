package mysql

import (
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/ipam"
)

// CreateSubnet create a subnet
func (m *MySQL) CreateSubnet(ctx context.Context, subnet ipam.Subnet) (*ipam.Subnet, error) {
	query := `INSERT INTO subnet(uuid, name, network, start, end, gateway, dns_server, metadata_server) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, subnet.UUID, subnet.Name, subnet.Network, subnet.Start, subnet.End, subnet.Gateway, subnet.DNSServer, subnet.MetadataServer)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &subnet, nil
}

// GetSubnetByID retrieves address according to the id given
func (m *MySQL) GetSubnetByID(ctx context.Context, uuid uuid.UUID) (*ipam.Subnet, error) {
	query := `SELECT uuid, name, network, start, end, gateway, dns_server, metadata_server, created_at, updated_at FROM subnet WHERE uuid = ?`
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
