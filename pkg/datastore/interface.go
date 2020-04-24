package datastore

import "context"

type Datastore interface {
	GetIQN(ctx context.Context, hostname string) (string, error)
}
