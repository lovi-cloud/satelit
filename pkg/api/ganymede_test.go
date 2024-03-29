package api

import (
	"database/sql"
	"testing"

	"github.com/lovi-cloud/satelit/pkg/ganymede"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-test/deep"
	pb "github.com/lovi-cloud/satelit/api/satelit"
)

func TestSatelitServer_AddVirtualMachine(t *testing.T) {
	nodes := getDummyNodes()
	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	cpgResp, err := client.AddCPUPinningGroup(ctx, &pb.AddCPUPinningGroupRequest{
		Name:           "testgroup",
		CountOfCore:    2,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to add cpu pinning group: %+v", err)
	}

	tests := []struct {
		input *pb.AddVirtualMachineRequest
		want  *pb.AddVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.AddVirtualMachineRequest{
				Name:              "test001",
				Vcpus:             1,
				MemoryKib:         1 * 1024 * 1024,
				RootVolumeGb:      10,
				SourceImageId:     imageResp.Image.Id,
				HypervisorName:    hypervisorName,
				EuropaBackendName: testEuropaBackendName,
			},
			want: &pb.AddVirtualMachineResponse{
				Uuid: "",
				Name: "test001",
			},
			err: false,
		},
		{
			input: &pb.AddVirtualMachineRequest{
				Name:              "test-with-cpu-pinning-group",
				Vcpus:             1,
				MemoryKib:         1 * 1024 * 1024,
				RootVolumeGb:      10,
				SourceImageId:     imageResp.Image.Id,
				HypervisorName:    hypervisorName,
				PinningGroupName:  cpgResp.CpuPinningGroup.Name,
				EuropaBackendName: testEuropaBackendName,
			},
			want: &pb.AddVirtualMachineResponse{
				Uuid: "",
				Name: "test-with-cpu-pinning-group",
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.AddVirtualMachine(ctx, test.input)
		if got != nil {
			test.want.Uuid = got.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_StartVirtualMachine(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}

	tests := []struct {
		input *pb.StartVirtualMachineRequest
		want  *pb.StartVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.StartVirtualMachineRequest{
				Uuid: vmResp.Uuid,
			},
			want: &pb.StartVirtualMachineResponse{
				Uuid: vmResp.Uuid,
				Name: vmResp.Uuid, // TODO
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.StartVirtualMachine(ctx, test.input)
		if got != nil {
			test.want.Uuid = got.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ShowVirtualMachine(t *testing.T) {
	nodes := getDummyNodes()
	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	cpgResp, err := client.AddCPUPinningGroup(ctx, &pb.AddCPUPinningGroupRequest{
		Name:           "testgroup",
		CountOfCore:    2,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to add cpu pinning group: %+v", err)
	}

	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}

	vmCpgResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test-with-cpu-pinning-group",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		PinningGroupName:  cpgResp.CpuPinningGroup.Name,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine with cpu pinning group: %+v", err)
	}

	tests := []struct {
		input *pb.ShowVirtualMachineRequest
		want  *pb.ShowVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.ShowVirtualMachineRequest{Uuid: vmResp.Uuid},
			want: &pb.ShowVirtualMachineResponse{
				VirtualMachine: &pb.VirtualMachine{
					Uuid:              vmResp.Uuid,
					Name:              "test001",
					Vcpus:             1,
					MemoryKib:         1 * 1024 * 1024,
					HypervisorName:    hypervisorName,
					SourceImageId:     imageResp.Image.Id,
					RootVolumeGb:      10,
					EuropaBackendName: testEuropaBackendName,
				},
			},
			err: false,
		},
		{
			input: &pb.ShowVirtualMachineRequest{Uuid: vmCpgResp.Uuid},
			want: &pb.ShowVirtualMachineResponse{
				VirtualMachine: &pb.VirtualMachine{
					Uuid:              vmCpgResp.Uuid,
					Name:              "test-with-cpu-pinning-group",
					Vcpus:             1,
					MemoryKib:         1 * 1024 * 1024,
					RootVolumeGb:      10,
					SourceImageId:     imageResp.Image.Id,
					HypervisorName:    hypervisorName,
					PinningGroupName:  cpgResp.CpuPinningGroup.Name,
					EuropaBackendName: testEuropaBackendName,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ShowVirtualMachine(ctx, test.input)
		if got != nil {
			test.want.VirtualMachine.Uuid = got.VirtualMachine.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListVirtualMachine(t *testing.T) {
	nodes := getDummyNodes()
	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	cpgResp, err := client.AddCPUPinningGroup(ctx, &pb.AddCPUPinningGroupRequest{
		Name:           "testgroup",
		CountOfCore:    2,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to add cpu pinning group: %+v", err)
	}

	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}

	vmCpgResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test-with-cpu-pinning-group",
		Vcpus:             2, // sort key for test
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		PinningGroupName:  cpgResp.CpuPinningGroup.Name,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine with cpu pinning group: %+v", err)
	}

	tests := []struct {
		input *pb.ListVirtualMachineRequest
		want  *pb.ListVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.ListVirtualMachineRequest{},
			want: &pb.ListVirtualMachineResponse{
				VirtualMachines: []*pb.VirtualMachine{
					{
						Uuid:              vmResp.Uuid,
						Name:              "test001",
						Vcpus:             1,
						MemoryKib:         1 * 1024 * 1024,
						RootVolumeGb:      10,
						SourceImageId:     imageResp.Image.Id,
						HypervisorName:    hypervisorName,
						EuropaBackendName: testEuropaBackendName,
					},
					{
						Uuid:              vmCpgResp.Uuid,
						Name:              "test-with-cpu-pinning-group",
						Vcpus:             2,
						MemoryKib:         1 * 1024 * 1024,
						RootVolumeGb:      10,
						SourceImageId:     imageResp.Image.Id,
						HypervisorName:    hypervisorName,
						PinningGroupName:  cpgResp.CpuPinningGroup.Name,
						EuropaBackendName: testEuropaBackendName,
					},
				},
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := client.ListVirtualMachine(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteVirtualMachine(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}

	tests := []struct {
		input *pb.DeleteVirtualMachineRequest
		want  *pb.DeleteVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.DeleteVirtualMachineRequest{
				Uuid: vmResp.Uuid,
			},
			want: &pb.DeleteVirtualMachineResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteVirtualMachine(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_CreateBridge(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	tests := []struct {
		input *pb.CreateBridgeRequest
		want  *pb.CreateBridgeResponse
		err   bool
	}{
		{
			input: &pb.CreateBridgeRequest{
				Name:   "testbr1000",
				VlanId: 1000,
			},
			want: &pb.CreateBridgeResponse{Bridge: &pb.Bridge{
				Uuid:   "",
				VlanId: 1000,
				Name:   "testbr1000",
			}},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.CreateBridge(ctx, test.input)
		if got != nil {
			test.want.Bridge.Uuid = got.Bridge.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_CreateInternalBridge(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	tests := []struct {
		input *pb.CreateInternalBridgeRequest
		want  *pb.CreateInternalBridgeResponse
		err   bool
	}{
		{
			input: &pb.CreateInternalBridgeRequest{
				Name: "testbr1000",
			},
			want: &pb.CreateInternalBridgeResponse{Bridge: &pb.Bridge{
				Uuid:   "",
				VlanId: 0,
				Name:   "testbr1000",
			}},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.CreateInternalBridge(ctx, test.input)
		if got != nil {
			test.want.Bridge.Uuid = got.Bridge.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_GetBridge(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		input *pb.GetBridgeRequest
		want  *pb.GetBridgeResponse
		err   bool
	}{
		{
			input: &pb.GetBridgeRequest{
				Uuid: resp.Bridge.Uuid,
			},
			want: &pb.GetBridgeResponse{
				Bridge: resp.Bridge,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.GetBridge(ctx, test.input)
		if got != nil {
			test.want.Bridge.Uuid = got.Bridge.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListBridge(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		input *pb.ListBridgeRequest
		want  *pb.ListBridgeResponse
		err   bool
	}{
		{
			input: &pb.ListBridgeRequest{},
			want: &pb.ListBridgeResponse{
				Bridges: []*pb.Bridge{
					resp.Bridge,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListBridge(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteBridge(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	tests := []struct {
		input *pb.DeleteBridgeRequest
		want  *pb.DeleteBridgeResponse
		err   bool
	}{
		{
			input: &pb.DeleteBridgeRequest{
				Uuid: resp.Bridge.Uuid,
			},
			want: &pb.DeleteBridgeResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteBridge(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_AttachInterface(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}
	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}
	bridgeResp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create test bridge: %+v", err)
	}
	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "testsubnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input *pb.AttachInterfaceRequest
		want  *pb.AttachInterfaceResponse
		err   bool
	}{
		{
			input: &pb.AttachInterfaceRequest{
				VirtualMachineId: vmResp.Uuid,
				BridgeId:         bridgeResp.Bridge.Uuid,
				Average:          1 * 1024 * 1024,
				Name:             "vnet0",
				LeaseId:          leaseResp.Lease.Uuid,
			},
			want: &pb.AttachInterfaceResponse{
				InterfaceAttachment: &pb.InterfaceAttachment{
					Uuid:             "",
					VirtualMachineId: vmResp.Uuid,
					BridgeId:         bridgeResp.Bridge.Uuid,
					Average:          1 * 1024 * 1024,
					Name:             "vnet0",
					LeaseId:          leaseResp.Lease.Uuid,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.AttachInterface(ctx, test.input)
		if got != nil {
			test.want.InterfaceAttachment.Uuid = got.InterfaceAttachment.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DetachInterface(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}
	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}
	bridgeResp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create test bridge: %+v", err)
	}
	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "testsubnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}
	attachResp, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: vmResp.Uuid,
		BridgeId:         bridgeResp.Bridge.Uuid,
		Average:          1 * 1024 * 1024,
		Name:             "vnet0",
		LeaseId:          leaseResp.Lease.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to attach test interface: %+v", err)
	}

	tests := []struct {
		input *pb.DetachInterfaceRequest
		want  *pb.DetachInterfaceResponse
		err   bool
	}{
		{
			input: &pb.DetachInterfaceRequest{
				AtttachmentId: attachResp.InterfaceAttachment.Uuid,
			},
			want: &pb.DetachInterfaceResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DetachInterface(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_GetAttachment(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}
	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}
	bridgeResp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create test bridge: %+v", err)
	}
	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "testsubnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}
	attachResp, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: vmResp.Uuid,
		BridgeId:         bridgeResp.Bridge.Uuid,
		Average:          1 * 1024 * 1024,
		Name:             "vnet0",
		LeaseId:          leaseResp.Lease.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to attach test interface: %+v", err)
	}

	tests := []struct {
		input *pb.GetAttachmentRequest
		want  *pb.GetAttachmentResponse
		err   bool
	}{
		{
			input: &pb.GetAttachmentRequest{
				AttachmentId: attachResp.InterfaceAttachment.Uuid,
			},
			want: &pb.GetAttachmentResponse{
				InterfaceAttachment: attachResp.InterfaceAttachment,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.GetAttachment(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListAttachment(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}
	vmResp, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              "test001",
		Vcpus:             1,
		MemoryKib:         1 * 1024 * 1024,
		RootVolumeGb:      10,
		SourceImageId:     imageResp.Image.Id,
		HypervisorName:    hypervisorName,
		EuropaBackendName: testEuropaBackendName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
	}
	bridgeResp, err := client.CreateBridge(ctx, &pb.CreateBridgeRequest{
		Name:   "testbr1000",
		VlanId: 1000,
	})
	if err != nil {
		t.Fatalf("failed to create test bridge: %+v", err)
	}
	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "testsubnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}
	attachResp, err := client.AttachInterface(ctx, &pb.AttachInterfaceRequest{
		VirtualMachineId: vmResp.Uuid,
		BridgeId:         bridgeResp.Bridge.Uuid,
		Average:          1 * 1024 * 1024,
		Name:             "vnet0",
		LeaseId:          leaseResp.Lease.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to attach test interface: %+v", err)
	}

	tests := []struct {
		input *pb.ListAttachmentRequest
		want  *pb.ListAttachmentResponse
		err   bool
	}{
		{
			input: &pb.ListAttachmentRequest{},
			want: &pb.ListAttachmentResponse{
				InterfaceAttachments: []*pb.InterfaceAttachment{
					attachResp.InterfaceAttachment,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListAttachment(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_AddCPUPinningGroup(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	nodes := getDummyNodes()

	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	tests := []struct {
		input   *pb.AddCPUPinningGroupRequest
		want    *pb.AddCPUPinningGroupResponse
		errCode codes.Code
	}{
		{
			input: &pb.AddCPUPinningGroupRequest{
				Name:           "testgroup",
				CountOfCore:    2,
				HypervisorName: hypervisorName,
			},
			want: &pb.AddCPUPinningGroupResponse{CpuPinningGroup: &pb.CPUPinningGroup{
				Uuid:        "",
				Name:        "testgroup",
				CountOfCore: 2,
			}},
			errCode: codes.OK,
		},
		{
			input: &pb.AddCPUPinningGroupRequest{
				Name:           "not_multiple_of_two_group",
				CountOfCore:    3,
				HypervisorName: hypervisorName,
			},
			want:    nil,
			errCode: codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		got, err := client.AddCPUPinningGroup(ctx, test.input)
		if got != nil {
			test.want.CpuPinningGroup.Uuid = got.CpuPinningGroup.Uuid
		}
		if test.errCode == 0 && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}

		s, ok := status.FromError(err)
		if test.errCode != 0 && ok && s.Code() != test.errCode {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ShowCPUPinningGroup(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	nodes := getDummyNodes()

	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	resp, err := client.AddCPUPinningGroup(ctx, &pb.AddCPUPinningGroupRequest{
		Name:           "testgroup",
		CountOfCore:    4,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to addCPUPinningGroup: %+v", err)
	}

	tests := []struct {
		input   *pb.ShowCPUPinningGroupRequest
		want    *pb.ShowCPUPinningGroupResponse
		errCode codes.Code
	}{
		{
			input: &pb.ShowCPUPinningGroupRequest{
				Uuid: resp.CpuPinningGroup.Uuid,
			},
			want: &pb.ShowCPUPinningGroupResponse{
				CpuPinningGroup: &pb.CPUPinningGroup{
					Uuid:        resp.CpuPinningGroup.Uuid,
					Name:        "testgroup",
					CountOfCore: 4,
				},
			},
			errCode: codes.OK,
		},
	}

	for _, test := range tests {
		got, err := client.ShowCPUPinningGroup(ctx, test.input)
		if test.errCode == 0 && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}

		s, ok := status.FromError(err)
		if test.errCode != 0 && ok && s.Code() != test.errCode {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteCPUPinningGroup(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	nodes := getDummyNodes()

	hypervisorName, teardownTeleskop, err := setupTeleskop(nodes)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	resp, err := client.AddCPUPinningGroup(ctx, &pb.AddCPUPinningGroupRequest{
		Name:           "testgroup",
		CountOfCore:    4,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to addCPUPinningGroup: %+v", err)
	}

	tests := []struct {
		input   *pb.DeleteCPUPinningGroupRequest
		want    *pb.DeleteCPUPinningGroupResponse
		errCode codes.Code
	}{
		{
			input: &pb.DeleteCPUPinningGroupRequest{
				Uuid: resp.CpuPinningGroup.Uuid,
			},
			want:    &pb.DeleteCPUPinningGroupResponse{},
			errCode: codes.OK,
		},
	}

	for _, test := range tests {
		got, err := client.DeleteCPUPinningGroup(ctx, test.input)
		if test.errCode == 0 && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}

		s, ok := status.FromError(err)
		if test.errCode != 0 && ok && s.Code() != test.errCode {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func getDummyNodes() []ganymede.NUMANode {
	var (
		testNUMANodeUUID  = "162b42f5-2eea-4fd1-b57b-c598db69fb4a"
		testCorePairUUIDs = []uuid.UUID{
			uuid.FromStringOrNil("9cf11645-ec85-4607-b638-cd592819bbae"),
			uuid.FromStringOrNil("25b403a9-cdd7-4176-8d44-c922220bdcb8"),
			uuid.FromStringOrNil("2cc61359-8912-4187-aadc-8692574b1b52"),
			uuid.FromStringOrNil("e77523a3-fef0-4864-b24f-4f9579a65eed"),
		}
	)

	var cps []ganymede.CorePair
	for i, u := range testCorePairUUIDs {

		l := int32(len(testCorePairUUIDs) + i)
		ll := sql.NullInt32{
			Int32: l,
			Valid: true,
		}

		cp := ganymede.CorePair{
			UUID:         u,
			PhysicalCore: uint32(i),
			LogicalCore:  ll,
			NUMANodeID:   uuid.FromStringOrNil(testNUMANodeUUID),
		}
		cps = append(cps, cp)
	}

	nodes := []ganymede.NUMANode{
		{
			UUID:            uuid.FromStringOrNil(testNUMANodeUUID),
			CorePairs:       cps,
			PhysicalCoreMin: 0,
			PhysicalCoreMax: 3,
			LogicalCoreMin: sql.NullInt32{
				Int32: 4,
				Valid: true,
			},
			LogicalCoreMax: sql.NullInt32{
				Int32: 8,
				Valid: true,
			},
		},
	}

	return nodes
}
