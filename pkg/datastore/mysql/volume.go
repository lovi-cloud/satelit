package mysql

import (
	"fmt"

	"github.com/whywaita/satelit/pkg/europa"
)

// GetVolume return volume from datastore
func (m *MySQL) GetVolume(volumeID string) (*europa.Volume, error) {
	var v europa.Volume

	query := fmt.Sprintf(`SELECT * FROM volume WHERE id = '%s'`, volumeID)
	err := m.Conn.Get(&v, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &v, nil
}

// PutVolume write volume record.
func (m *MySQL) PutVolume(volume europa.Volume) error {
	query := `INSERT INTO volume(id, attached, hostname, capacity_gb, base_image_id, host_lun_id) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.Conn.Exec(query, volume.ID, volume.Attached, volume.HostName, volume.CapacityGB, volume.BaseImageID, volume.HostLUNID)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
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