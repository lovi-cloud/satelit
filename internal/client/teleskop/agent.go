package teleskop

import (
	"fmt"
	"sync"

	agentpb "github.com/whywaita/satelit/api"
	"google.golang.org/grpc"
)

var (
	client map[string]agentpb.AgentClient
	mu     sync.RWMutex
)

// New create teleskop map
func New(endpoints map[string]string) error {
	c := make(map[string]agentpb.AgentClient)

	for hostname, endpoint := range endpoints {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithInsecure(),
		)
		if err != nil {
			return fmt.Errorf("failed to connect teleskop endpoint: %w", err)
		}

		mu.Lock()
		c[hostname] = agentpb.NewAgentClient(conn)
		mu.Unlock()
	}

	client = c
	return nil
}

// GetClient return teleskop Client
func GetClient(hostname string) agentpb.AgentClient {
	mu.RLock()
	c := client[hostname]
	mu.RUnlock()

	return c
}
