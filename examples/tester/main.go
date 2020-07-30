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

	resp11, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "dummy",
		VlanId:         0,
		Network:        "172.19.0.0/16",
		Start:          "172.19.0.2",
		End:            "172.19.255.254",
		Gateway:        "172.19.0.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "172.19.0.1",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp11)

	resp7, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "v1002",
		VlanId:         1002,
		Network:        "10.160.2.0/24",
		Start:          "10.160.2.100",
		End:            "10.160.2.200",
		Gateway:        "10.160.2.254",
		DnsServer:      "8.8.8.8",
		MetadataServer: "10.160.2.15",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp7)

	resp1, err := client.CreateBridge(context.Background(), &pb.CreateBridgeRequest{
		Name:   "br1002",
		VlanId: 1002,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp1)

	resp9, err := client.CreateInternalBridge(context.Background(), &pb.CreateInternalBridgeRequest{
		Name: "br-team001",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp9)

	resp3, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:           "yjuba-test006",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   40,
		SourceImageId:  "111abd25-224a-414f-a4ee-845f3e0845bb",
		HypervisorName: "whywaita-amd001",
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp3)

	resp6, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: resp7.Subnet.Uuid,
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
	fmt.Printf("%+v\n", resp5)

	resp12, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: resp11.Subnet.Uuid,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp12)

	resp13, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: resp12.Address.Uuid,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp13)

	resp4, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: resp3.Uuid,
		BridgeId:         resp1.Bridge.Uuid,
		Average:          1 * 1024 * 1024,
		Name:             "yjuba-net007",
		LeaseId:          resp5.Lease.Uuid,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp4)

	resp10, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: resp3.Uuid,
		BridgeId:         resp9.Bridge.Uuid,
		Average:          1 * 1024 * 1024,
		Name:             "yjuba-net008",
		LeaseId:          resp13.Lease.Uuid,
	})
	fmt.Printf("%+v\n", resp10)

	resp8, err := client.StartVirtualMachine(ctx, &pb.StartVirtualMachineRequest{
		Uuid: resp3.Uuid,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp8)

	return nil
}
