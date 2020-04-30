package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/whywaita/satelit/api/satelit"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:9262", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := pb.NewSatelitClient(conn)

	resp, err := client.GetVolumes(ctx, &pb.GetVolumesRequest{})
	if err != nil {
		return err
	}

	for _, v := range resp.Volumes {
		fmt.Printf("%+v\n", v)
	}

	return nil
}
