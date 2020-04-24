package teleskop

import (
	"sync"

	"github.com/pkg/errors"
	pb "github.com/whywaita/satelit/api"
	"google.golang.org/grpc"
)

var (
	client map[string]pb.AgentClient
	mu     sync.RWMutex
)

func New(endpoints []string) error {
	c := make(map[string]pb.AgentClient)

	for _, endpoint := range endpoints {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithInsecure(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to connect teleskop endpoint")
		}

		mu.Lock()
		c[endpoint] = pb.NewAgentClient(conn)
		mu.Unlock()
	}

	client = c
	return nil
}

func GetClient(endpoint string) pb.AgentClient {
	mu.RLock()
	c := client[endpoint]
	mu.RUnlock()

	return c
}
