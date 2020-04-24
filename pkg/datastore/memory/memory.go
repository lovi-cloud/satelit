package memory

import "context"

type Memory struct{}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) GetIQN(ctx context.Context, hostname string) (string, error) {
	return "dummy iqn", nil
}
