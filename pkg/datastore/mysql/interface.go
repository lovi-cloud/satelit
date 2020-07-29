package mysql

import (
	"context"
	"fmt"

	"github.com/whywaita/satelit/pkg/ganymede"
)

// AttachInterface is
func (m *MySQL) AttachInterface(ctx context.Context, attachment ganymede.InterfaceAttachment) (*ganymede.InterfaceAttachment, error) {
	query := `INSERT INTO interface_attachment(virtual_machine_id, bridge_id, average, name, lease_id) VALUES (?, ?, ?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	r, err := stmt.ExecContext(ctx, attachment.VirtualMachineID, attachment.BridgeID, attachment.Average, attachment.Name, attachment.LeaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert attachment: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return m.GetAttachment(ctx, int(id))
}

// DetachInterface is
func (m *MySQL) DetachInterface(ctx context.Context, attachmentID int) error {
	query := `DELETE FROM interface_attachemnt from id = ?`
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
func (m *MySQL) GetAttachment(ctx context.Context, attachmentID int) (*ganymede.InterfaceAttachment, error) {
	query := `SELECT id, virtual_machine_id, bridge_id, average, name, lease_id, created_at, updated_at FROM interface_attachment WHERE id = ?`
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
	query := `SELECT id, virtual_machine_id, bridge_id, average, name, lease_id, created_at, updated_at FROM interface_attachment`
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
