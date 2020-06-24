package mysql

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/pkg/europa"
)

// GetImage retrieves image object
func (m *MySQL) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	var image europa.BaseImage

	query := fmt.Sprintf(`SELECT * FROM image WHERE uuid = "%s"`, imageID)
	err := m.Conn.Get(&image, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &image, nil
}

// ListImage retrieves all images
func (m *MySQL) ListImage() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	query := fmt.Sprintf("SELECT * FROM image")
	err := m.Conn.Select(&images, query)
	if err != nil {
		return nil, fmt.Errorf("failed to SELCT image table: %w", err)
	}

	return images, nil
}

// PutImage write image record
// need to fix if call more than once
func (m *MySQL) PutImage(image europa.BaseImage) error {
	query := `INSERT INTO image(uuid, name, volume_id, description) VALUES (UUID_TO_BIN(?), ?, ?, ?)`
	_, err := m.Conn.Exec(query, image.UUID, image.Name, image.CacheVolumeID, image.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// DeleteImage delete image record
func (m *MySQL) DeleteImage(imageID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM image WHERE uuid = "%s"`, imageID)
	_, err := m.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
