package main

import (
	"context"
	"fmt"
	"os"

	"github.com/whywaita/satelit/internal/config"

	uuid "github.com/satori/go.uuid"

	pb "github.com/whywaita/satelit/api/satelit"
	"google.golang.org/grpc"
)

func init() {
	conf := "./../../configs/satelit.yaml"
	err := config.Load(&conf)
	if err != nil {
		panic(err)
	}
}

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

	fmt.Println("GetVolumes")
	_, err = client.GetVolumes(ctx, &pb.GetVolumesRequest{})
	if err != nil {
		return err
	}

	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	fmt.Println("AddVolume")
	resp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Id:           u.String(),
		CapacityByte: 23,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp.Volume)

	var hostname string
	for h := range config.GetValue().Teleskop.Endpoints {
		hostname = h
	}

	fmt.Println("AttachVolume")
	_, err = client.AttachVolume(ctx, &pb.AttachVolumeRequest{
		Id:       u.String(),
		Hostname: hostname,
	})
	if err != nil {
		return err
	}

	return nil
}
