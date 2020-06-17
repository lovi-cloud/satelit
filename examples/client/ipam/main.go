package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/config"
	"google.golang.org/grpc"
)

func init() {
	conf := "./../../../configs/satelit.yaml"
	err := config.Load(&conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
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

	fmt.Println("ListSubnet")
	resp1, err := client.ListSubnet(ctx, &pb.ListSubnetRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("%#q\n", resp1)

	fmt.Println("CreateSubnet")
	resp2, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "yjuba-test001",
		Network:        "10.192.0.0/23",
		Start:          "10.192.0.100",
		End:            "10.192.0.200",
		Gateway:        "10.192.0.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "10.192.0.15",
	})
	if err != nil {
		return err
	}

	fmt.Printf("%#q\n", resp2)

	fmt.Println("ListSubnet")
	resp3, err := client.ListSubnet(ctx, &pb.ListSubnetRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("%#q\n", resp3)

	return nil
}
