package cmd

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/lovi-cloud/satelit/api/satelit"
	"github.com/spf13/cobra"
)

var (
	vmName         string
	imageUUID      string
	vcpu           uint32
	memory         uint64
	rootVolumeGB   uint32
	hypervisorName string
)

var createVMCmd = &cobra.Command{
	Use:   "create-vm",
	Short: "Create and Start virtual machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(satelitAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			return err
		}

		client := pb.NewSatelitClient(conn)
		return SampleStartVirtualMachine(context.Background(), client)
	},
}

func init() {
	createVMCmd.Flags().StringVarP(&vmName, "name", "", "", "name of virtual machine")
	createVMCmd.MarkFlagRequired("name")
	createVMCmd.Flags().StringVarP(&imageUUID, "image", "", "", "uuid of image")
	createVMCmd.MarkFlagRequired("image")
	createVMCmd.Flags().Uint32VarP(&vcpu, "vcpu", "", 1, "number of cpu")
	createVMCmd.Flags().Uint64VarP(&memory, "memory", "", 2*1024*1024, "kib size of memory")
	createVMCmd.Flags().Uint32VarP(&rootVolumeGB, "gb", "", 21, "gb size of root disk")

	createVMCmd.Flags().StringVarP(&europaBackendName, "backend", "", "", "backend name of europa")
	createVMCmd.MarkFlagRequired("backend")
	createVMCmd.Flags().StringVarP(&hypervisorName, "hypervisor", "", "", "name of hypervisor")
	createVMCmd.MarkFlagRequired("hypervisor")

	rootCmd.AddCommand(createVMCmd)
}

// SampleStartVirtualMachine is sample of AddVirtualMachine / StartVirtualMachine
func SampleStartVirtualMachine(ctx context.Context, client pb.SatelitClient) error {
	fmt.Println("AddVirtualMachine")
	resp1, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              vmName,
		Vcpus:             vcpu,
		MemoryKib:         memory,
		RootVolumeGb:      rootVolumeGB,
		SourceImageId:     imageUUID,
		HypervisorName:    hypervisorName,
		EuropaBackendName: europaBackendName,
	})
	if err != nil {
		return err
	}
	vmUUID := resp1.Uuid

	fmt.Println("StartVirtualMachine")
	resp2, err := client.StartVirtualMachine(ctx, &pb.StartVirtualMachineRequest{Uuid: vmUUID})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp2)

	return nil
}
