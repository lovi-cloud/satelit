package teleskop

import (
	"sync"

	"github.com/pkg/errors"
	agentpb "github.com/whywaita/satelit/api"
	"google.golang.org/grpc"
)

var (
	client map[string]agentpb.AgentClient
	mu     sync.RWMutex
)

func New(endpoints map[string]string) error {
	c := make(map[string]agentpb.AgentClient)

	for hostname, endpoint := range endpoints {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithInsecure(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to connect teleskop endpoint")
		}

		mu.Lock()
		c[hostname] = agentpb.NewAgentClient(conn)
		mu.Unlock()
	}

	client = c
	return nil
}

func GetClient(hostname string) agentpb.AgentClient {
	mu.RLock()
	c := client[hostname]
	mu.RUnlock()

	return c
}
