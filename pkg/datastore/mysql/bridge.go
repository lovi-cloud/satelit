package mysql

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/ganymede"
)

// CreateBridge is
func (m *MySQL) CreateBridge(ctx context.Context, bridge ganymede.Bridge) (*ganymede.Bridge, error) {
	query := `INSERT INTO bridge(uuid, vlan_id, name) VALUES(?, ?, ?)`
	_, err := m.Conn.ExecContext(ctx, query, bridge.UUID, bridge.VLANID, bridge.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return &bridge, nil
}

// GetBridge is
func (m *MySQL) GetBridge(ctx context.Context, bridgeID uuid.UUID) (*ganymede.Bridge, error) {
	query := `SELECT uuid, vlan_id, name, created_at, updated_at FROM bridge where id = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	var bridge ganymede.Bridge
	err = stmt.GetContext(ctx, &bridge, bridgeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bridge: %w", err)
	}

	return &bridge, nil
}

// ListBridge is
func (m *MySQL) ListBridge(ctx context.Context) ([]ganymede.Bridge, error) {
	query := `SELECT uuid, vlan_id, name, created_at, updated_at FROM bridge`
	var bridges []ganymede.Bridge
	err := m.Conn.SelectContext(ctx, &bridges, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get bridge list: %w", err)
	}

	return bridges, nil
}

// DeleteBridge is
func (m *MySQL) DeleteBridge(ctx context.Context, bridgeID uuid.UUID) error {
	query := `DELETE FROM bridge WHERE id = ?`
	stmt, err := m.Conn.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	_, err = stmt.ExecContext(ctx, bridgeID)
	if err != nil {
		return fmt.Errorf("failed to delete bridge: %w", err)
	}

	return nil
}
