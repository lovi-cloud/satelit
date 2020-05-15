package memory

import "context"

// A Memory is on memory datastore for testing.
type Memory struct{}

// New create Memory
func New() *Memory {
	return &Memory{}
}

// GetIQN return IQN from on memory
func (m *Memory) GetIQN(ctx context.Context, hostname string) (string, error) {
	return "dummy iqn", nil
}
