package mysql

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/whywaita/satelit/pkg/europa"
)

// ListVolume retrieves multi volumes from MySQL IN query
func (m *MySQL) ListVolume(volumeIDs []string) ([]europa.Volume, error) {
	q, a, err := sqlx.In(`SELECT id, attached, hostname, capacity_gb, BIN_TO_UUID(base_image_id) as base_image_id, host_lun_id FROM volume WHERE id IN (?)`, volumeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	var volumes []europa.Volume
	err = m.Conn.Select(&volumes, q, a...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieves volumes: %w", err)
	}

	return volumes, nil
}

// GetVolume return volume from datastore
func (m *MySQL) GetVolume(volumeID string) (*europa.Volume, error) {
	var v europa.Volume

	query := fmt.Sprintf(`SELECT id, attached, hostname, capacity_gb, BIN_TO_UUID(base_image_id) as base_image_id, host_lun_id FROM volume WHERE id = '%s'`, volumeID)
	err := m.Conn.Get(&v, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &v, nil
}

// PutVolume write volume record.
func (m *MySQL) PutVolume(volume europa.Volume) error {
	_, err := m.GetVolume(volume.ID)
	if errors.Is(err, sql.ErrNoRows) {
		// no rows, need to insert
		query := `INSERT INTO volume(id, attached, hostname, capacity_gb, base_image_id, host_lun_id) VALUES (?, ?, ?, ?, UUID_TO_BIN(?), ?)`
		_, err := m.Conn.Exec(query, volume.ID, volume.Attached, volume.HostName, volume.CapacityGB, volume.BaseImageID, volume.HostLUNID)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to execute get query: %w", err)
	}

	// found rows, need to update
	query := `UPDATE volume SET attached=:attached, hostname=:hostname, capacity_gb=:capacityGB, host_lun_id=:hostLUNID WHERE id = :id`
	_, err = m.Conn.NamedExec(query, map[string]interface{}{
		"attached":   volume.Attached,
		"hostname":   volume.HostName,
		"hostLUNID":  volume.HostLUNID,
		"capacityGB": volume.CapacityGB,
		"id":         volume.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

// DeleteVolume delete volume record
func (m *MySQL) DeleteVolume(volumeID string) error {
	query := fmt.Sprintf(`DELETE FROM volume WHERE id = '%s'`, volumeID)
	_, err := m.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
