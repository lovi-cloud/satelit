package mysql

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/whywaita/satelit/pkg/ganymede"
)

// GetHypervisor retrieve hypervisor
func (m *MySQL) GetHypervisor(ctx context.Context, hypervisorID int) (*ganymede.HyperVisor, error) {
	var hypervisor ganymede.HyperVisor

	query := `SELECT id, iqn, hostname, created_at, updated_at FROM hypervisor WHERE id = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Get(&hypervisor, hypervisorID); err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &hypervisor, nil
}

// GetHypervisorByHostname retrieve hypervisor by hostname
func (m *MySQL) GetHypervisorByHostname(ctx context.Context, hostname string) (*ganymede.HyperVisor, error) {
	var hypervisor ganymede.HyperVisor

	query := `SELECT id, iqn, hostname, created_at, updated_at FROM hypervisor WHERE hostname = ?`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	if err := stmt.Get(&hypervisor, hostname); err != nil {
		return nil, fmt.Errorf("failed to execute get query: %w", err)
	}

	return &hypervisor, nil
}

// PutHypervisor register hypervisor
// return hypervisor ID by MySQL AUTO INCREMENT.
func (m *MySQL) PutHypervisor(ctx context.Context, iqn, hostname string) (int, error) {
	query := `INSERT INTO hypervisor(iqn, hostname) VALUES (?, ?)`
	stmt, err := m.Conn.Preparex(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}

	result, err := stmt.ExecContext(ctx, iqn, hostname)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				// duplicate error, It's already exits
				hv, err := m.GetHypervisorByHostname(ctx, hostname)
				if err != nil {
					return 0, fmt.Errorf("failed to get hypervisor by hostname: %w", err)
				}
				return hv.ID, nil
			}
		}

		return 0, fmt.Errorf("failed to execute insert query: %w", err)
	}

	hypervisorID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return int(hypervisorID), nil
}

// PutHypervisorNUMANode register cpu cores by hypervisor
func (m *MySQL) PutHypervisorNUMANode(ctx context.Context, nodes []ganymede.NUMANode, hypervisorID int) error {
	for _, node := range nodes {
		nodeQuery := `INSERT INTO hypervisor_numa_node(uuid, physical_core_min, physical_core_max, logical_core_min, logical_core_max, hypervisor_id) VALUES (?, ?, ?, ?, ?, ?)`
		nodeStmt, err := m.Conn.Preparex(nodeQuery)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		_, err = nodeStmt.ExecContext(ctx, node.UUID, node.PhysicalCoreMin, node.PhysicalCoreMax, node.LogicalCoreMin, node.LogicalCoreMax, hypervisorID)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				if mysqlErr.Number == 1062 {
					// duplicate error, It's already exits
					// numa node is not change in same hypervisor ID. so this function is skip.
					return nil
				}
			}

			return fmt.Errorf("failed to execute insert query: %w", err)
		}

		for _, pair := range node.CorePairs {
			pairQuery := `INSERT INTO hypervisor_cpu_pair(uuid, numa_node_id, physical_core_number, logical_core_number) VALUES (?, ?, ?, ?)`
			pairStmt, err := m.Conn.Preparex(pairQuery)
			if err != nil {
				return fmt.Errorf("failed to prepare statement: %w", err)
			}
			if _, err := pairStmt.ExecContext(ctx, pair.UUID, node.UUID, pair.PhysicalCore, pair.LogicalCore); err != nil {
				return fmt.Errorf("failed to execute insert query: %w", err)
			}
		}
	}

	return nil
}
