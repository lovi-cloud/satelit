package mysql

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/lovi-cloud/satelit/pkg/ganymede"
)

// PutCPUPinningGroup put cpu pinning group
func (m *MySQL) PutCPUPinningGroup(ctx context.Context, cpuPinningGroup ganymede.CPUPinningGroup) error {
	query := `INSERT INTO cpu_pinning_group(uuid, name, hypervisor_id, count_of_core) VALUES (?, ?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	if _, err := stmt.ExecContext(ctx, cpuPinningGroup.UUID, cpuPinningGroup.Name, cpuPinningGroup.HypervisorID, cpuPinningGroup.CountCore); err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// GetCPUPinningGroup retrieve cpu pinning group
func (m *MySQL) GetCPUPinningGroup(ctx context.Context, cpuPinningGroupID uuid.UUID) (*ganymede.CPUPinningGroup, error) {
	query := `SELECT uuid, name, hypervisor_id, count_of_core, created_at, updated_at FROM cpu_pinning_group WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var cpg ganymede.CPUPinningGroup
	if err := stmt.GetContext(ctx, &cpg, cpuPinningGroupID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &cpg, nil
}

// GetCPUPinningGroupByName retrieve cpu pinning group by name
func (m *MySQL) GetCPUPinningGroupByName(ctx context.Context, name string) (*ganymede.CPUPinningGroup, error) {
	query := `SELECT uuid, name, hypervisor_id, count_of_core, created_at, updated_at FROM cpu_pinning_group WHERE name = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var cpg ganymede.CPUPinningGroup
	if err := stmt.GetContext(ctx, &cpg, name); err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &cpg, nil
}

// DeleteCPUPinningGroup delete cpu pinning group record
func (m *MySQL) DeleteCPUPinningGroup(ctx context.Context, cpuPinningGroupID uuid.UUID) error {
	query := `DELETE FROM cpu_pinning_group WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	if _, err := stmt.ExecContext(ctx, cpuPinningGroupID.String()); err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}

// GetAvailableCorePair retrieve not pinned CorePairs
func (m *MySQL) GetAvailableCorePair(ctx context.Context, hypervisorID int) ([]ganymede.NUMANode, error) {
	nodeQuery := `SELECT uuid FROM hypervisor_numa_node WHERE hypervisor_id = ?`
	nodeStmt, err := m.Conn.Preparex(nodeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var numaNodes []ganymede.NUMANode
	if err := nodeStmt.SelectContext(ctx, &numaNodes, hypervisorID); err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	// NOTE(whywaita): which one is best, OUTER JOIN or two queries? :thinking_face:
	pairQuery := `SELECT hcp.uuid, hcp.numa_node_id, hcp.physical_core_number, hcp.logical_core_number, hcp.created_at, hcp.updated_at FROM hypervisor_cpu_pair AS hcp
 LEFT OUTER JOIN cpu_core_pinned AS ccp ON (hcp.uuid = ccp.hypervisor_cpu_pair_id) WHERE ccp.uuid is NULL AND hcp.numa_node_id = ?`
	pairStmt, err := m.Conn.Preparex(pairQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	for i, node := range numaNodes {
		var p []ganymede.CorePair
		if err := pairStmt.Select(&p, node.UUID.String()); err != nil {
			return nil, fmt.Errorf("failed to execute select query: %w", err)
		}

		numaNodes[i].CorePairs = p
	}

	return numaNodes, nil
}

// GetCPUCorePair retrieve CPU CorePair
func (m *MySQL) GetCPUCorePair(ctx context.Context, corePairID uuid.UUID) (*ganymede.CorePair, error) {
	query := `SELECT uuid, numa_node_id, physical_core_number, logical_core_number, created_at, updated_at FROM hypervisor_cpu_pair WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var cp ganymede.CorePair
	if err := stmt.GetContext(ctx, &cp, corePairID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &cp, nil
}

// GetPinnedCoreByPinningGroup retrieve CPU CorePair
func (m *MySQL) GetPinnedCoreByPinningGroup(ctx context.Context, cpuPinningGroupID uuid.UUID) ([]ganymede.CPUCorePinned, error) {
	query := `SELECT uuid, pinning_group_id, hypervisor_cpu_pair_id, created_at, updated_at FROM cpu_core_pinned WHERE pinning_group_id = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var pinneds []ganymede.CPUCorePinned
	if err := stmt.SelectContext(ctx, &pinneds, cpuPinningGroupID.String()); err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return pinneds, nil
}

// PutPinnedCore put pinned record
func (m *MySQL) PutPinnedCore(ctx context.Context, pinned ganymede.CPUCorePinned) error {
	query := `INSERT INTO cpu_core_pinned(uuid, pinning_group_id, hypervisor_cpu_pair_id) VALUES (?, ?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	if _, err := stmt.ExecContext(ctx, pinned.UUID, pinned.CPUPinningGroupID, pinned.CorePairID); err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// DeletePinnedCore delete record
func (m *MySQL) DeletePinnedCore(ctx context.Context, pinnedID uuid.UUID) error {
	query := `DELETE FROM cpu_core_pinned WHERE uuid = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	if _, err := stmt.ExecContext(ctx, pinnedID); err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
