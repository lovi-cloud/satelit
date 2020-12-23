package mysql

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/pkg/ganymede"
)

// AttachInterface is
func (m *MySQL) AttachInterface(ctx context.Context, attachment ganymede.InterfaceAttachment) (*ganymede.InterfaceAttachment, error) {
	query := `INSERT INTO interface_attachment(uuid, virtual_machine_id, bridge_id, average, name, lease_id) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.ExecContext(ctx, attachment.UUID, attachment.VirtualMachineID, attachment.BridgeID, attachment.Average, attachment.Name, attachment.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert attachment: %w", err)
	}

	return m.GetAttachment(ctx, attachment.UUID)
}

// DetachInterface is
func (m *MySQL) DetachInterface(ctx context.Context, attachmentID uuid.UUID) error {
	query := `DELETE FROM interface_attachment WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.ExecContext(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	return nil
}

// GetAttachment is
func (m *MySQL) GetAttachment(ctx context.Context, attachmentID uuid.UUID) (*ganymede.InterfaceAttachment, error) {
	query := `SELECT uuid, virtual_machine_id, bridge_id, average, name, lease_id, created_at, updated_at FROM interface_attachment WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	var attachment ganymede.InterfaceAttachment
	err = stmt.GetContext(ctx, &attachment, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	return &attachment, nil
}

// ListAttachment is
func (m *MySQL) ListAttachment(ctx context.Context) ([]ganymede.InterfaceAttachment, error) {
	query := `SELECT uuid, virtual_machine_id, bridge_id, average, name, lease_id, created_at, updated_at FROM interface_attachment`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	var attachments []ganymede.InterfaceAttachment
	err = stmt.SelectContext(ctx, &attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment list: %w", err)
	}

	return attachments, nil
}
