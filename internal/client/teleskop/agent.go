package teleskop

import (
	"errors"
	"fmt"
	"sync"

	agentpb "github.com/whywaita/teleskop/protoc/agent"
	"google.golang.org/grpc"
)

var (
	client map[string]agentpb.AgentClient
	mu     sync.RWMutex
)

// Error const
var (
	ErrTeleskopAgentNotFound     = errors.New("a teleskop agent is not registered")
	ErrTeleskopAgentAlreadyExist = errors.New("a teleskop agent is already exist")
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
func GetClient(hostname string) (agentpb.AgentClient, error) {
	mu.RLock()
	c, ok := client[hostname]
	mu.RUnlock()

	if !ok {
		return nil, ErrTeleskopAgentNotFound
	}

	return c, nil
}

// ListClient return all tekeskop Clients
func ListClient() ([]agentpb.AgentClient, error) {
	mu.RLock()
	var cs []agentpb.AgentClient
	for _, c := range client {
		cs = append(cs, c)
	}
	mu.RUnlock()

	if len(cs) == 0 {
		return nil, ErrTeleskopAgentNotFound
	}

	return cs, nil
}

// AddClient add new teleskop Client and reconnect if registered
func AddClient(hostname, endpoint string) error {
	mu.Lock()
	defer mu.Unlock()

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to dial teleskop endpoint: %w", err)
	}
	client[hostname] = agentpb.NewAgentClient(conn)

	return nil
}
