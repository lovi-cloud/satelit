package mysql

import (
	"context"
	"fmt"

	"github.com/whywaita/satelit/pkg/ganymede"
)

// PutCPUPinningGroup put cpu pinning group
func (m *MySQL) PutCPUPinningGroup(ctx context.Context, cpuPinningGroup ganymede.CPUPinningGroup) error {
	query := `INSERT INTO cpu_pinning_group(uuid, name, count_of_core) VALUES (?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	if _, err := stmt.ExecContext(ctx, cpuPinningGroup.UUID, cpuPinningGroup.Name, cpuPinningGroup.CountCore); err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// GetAvailableCorePair retrieve not pinned CorePairs.
func (m *MySQL) GetAvailableCorePair(ctx context.Context, hypervisorID int) ([]ganymede.NUMANode, error) {
	panic("implement me")
}
