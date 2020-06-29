package mysql

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/europa"
)

// GetImage retrieves image object
func (m *MySQL) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	var image europa.BaseImage

	query := `SELECT uuid, name, description, volume_id, created_at, updated_at FROM image WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	err = stmt.Get(&image, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &image, nil
}

// ListImage retrieves all images
func (m *MySQL) ListImage() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	query := `SELECT uuid, name, description, volume_id, created_at, updated_at FROM image`
	stmt, err := m.Conn.Preparex(query)
	err = stmt.Select(&images)
	if err != nil {
		return nil, fmt.Errorf("failed to SELCT image table: %w", err)
	}

	return images, nil
}

// PutImage write image record
// need to fix if call more than once
func (m *MySQL) PutImage(image europa.BaseImage) error {
	query := `INSERT INTO image(uuid, name, volume_id, description) VALUES (?, ?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	_, err = stmt.Exec(image.UUID, image.Name, image.CacheVolumeID, image.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// DeleteImage delete image record
func (m *MySQL) DeleteImage(imageID uuid.UUID) error {
	query := `DELETE FROM image WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	_, err = stmt.Exec(imageID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
