package teleskop

import (
	"sync"

	"github.com/pkg/errors"
	pb "github.com/whywaita/satelit/api"
	"github.com/whywaita/satelit/internal/config"
	"google.golang.org/grpc"
)

var (
	client map[string]pb.AgentClient
	mu     sync.RWMutex
)

func New(endpoint string) error {
	c := make(map[string]pb.AgentClient)

	conn, err := grpc.Dial(
		config.GetValue().Teleskop.Endpoint,
		grpc.WithInsecure(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to connect teleskop endpoint")
	}

	mu.Lock()
	c[endpoint] = pb.NewAgentClient(conn)
	mu.Unlock()

	client = c
	return nil
}

func GetClient(endpoint string) pb.AgentClient {
	mu.RLock()
	c := client[endpoint]
	mu.RUnlock()

	return c
}
