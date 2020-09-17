package teleskop

import (
	"errors"
	"fmt"
	"sync"

	"github.com/whywaita/satelit/internal/logger"

	agentpb "github.com/whywaita/teleskop/protoc/agent"
	"google.golang.org/grpc"
)

var (
	connections map[string]*grpc.ClientConn
	mu          sync.RWMutex
)

// Error const
var (
	ErrTeleskopAgentNotFound     = errors.New("a teleskop agent is not registered")
	ErrTeleskopAgentAlreadyExist = errors.New("a teleskop agent is already exist")
)

// New create teleskop map
func New(endpoints map[string]string) error {
	c := make(map[string]*grpc.ClientConn)

	for hostname, endpoint := range endpoints {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithInsecure(),
		)
		if err != nil {
			return fmt.Errorf("failed to connect teleskop endpoint: %w", err)
		}

		mu.Lock()
		c[hostname] = conn
		mu.Unlock()
	}

	connections = c
	return nil
}

// GetClient return teleskop Client
func GetClient(hostname string) (agentpb.AgentClient, error) {
	mu.RLock()
	c, ok := connections[hostname]
	mu.RUnlock()

	if !ok {
		return nil, ErrTeleskopAgentNotFound
	}

	return agentpb.NewAgentClient(c), nil
}

// ListClient return all tekeskop Clients
func ListClient() ([]agentpb.AgentClient, error) {
	mu.RLock()
	var cs []agentpb.AgentClient
	for _, c := range connections {
		cs = append(cs, agentpb.NewAgentClient(c))
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

	conn, ok := connections[hostname]
	if ok {
		// already exits, close connection
		if err := conn.Close(); err != nil {
			logger.Logger.Debug(fmt.Sprintf("failed to close old teleskop connection: %+v", err))
		}
	}

	newConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to dial teleskop endpoint: %w", err)
	}
	connections[hostname] = newConn

	return nil
}
