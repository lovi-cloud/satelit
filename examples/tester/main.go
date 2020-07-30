package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/whywaita/satelit/api/satelit"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, "10.197.32.54:9262", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewSatelitClient(conn)

	resp2, err := client.ListBridge(ctx, &pb.ListBridgeRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("AAA: %+v\n", resp2)

	/*resp1, err := client.CreateBridge(context.Background(), &pb.CreateBridgeRequest{
		VlanId: 1001,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp1)*/

	/*resp3, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:           "yjuba-test003",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   40,
		SourceImageId:  "50cd4e98-ee92-475b-b991-a6204602ebb8",
		HypervisorName: "whywaita-amd001",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp3)*/

	/*resp7, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "v1001",
		Network:        "10.160.1.0/24",
		Start:          "10.160.1.100",
		End:            "10.160.1.200",
		Gateway:        "10.160.1.254",
		DnsServer:      "8.8.8.8",
		MetadataServer: "10.160.1.15",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp7)*/

	/*resp6, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: "e2a7d69a-8b80-4c90-9b00-714a41473a00",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp6)

	resp5, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: resp6.Address.Uuid,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp5)*/

	resp4, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: "85eb320e-5769-4566-87af-865adeb01be7",
		BridgeId:         "0e613f9e-25cb-4466-b7cb-88ec51e5df92",
		Average:          1 * 1024 * 1024,
		Name:             "yjuba-net003",
		LeaseId:          "b962c21f-eede-49e7-b01e-c2d83f165d13",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp4)

	return nil
}
