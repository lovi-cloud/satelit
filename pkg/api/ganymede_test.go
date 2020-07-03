package api

import (
	"testing"

	"github.com/go-test/deep"

	pb "github.com/whywaita/satelit/api/satelit"
)

func TestSatelitServer_AddVirtualMachine(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop()
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

	tests := []struct {
		input *pb.AddVirtualMachineRequest
		want  *pb.AddVirtualMachineResponse
		err   bool
	}{
		{
			input: &pb.AddVirtualMachineRequest{
				Name:           "test001",
				Vcpus:          1,
				MemoryKib:      1 * 1024 * 1024,
				RootVolumeGb:   10,
				SourceImageId:  imageResp.Image.Id,
				HypervisorName: hypervisorName,
			},
			want: &pb.AddVirtualMachineResponse{
				Uuid: "",
				Name: "test001",
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
	hypervisorName, teardownTeleskop, err := setupTeleskop()
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
		Name:           "test001",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   10,
		SourceImageId:  imageResp.Image.Id,
		HypervisorName: hypervisorName,
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
	hypervisorName, teardownTeleskop, err := setupTeleskop()
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
		Name:           "test001",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   10,
		SourceImageId:  imageResp.Image.Id,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to add test virtual machine: %+v\n", err)
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
					Uuid:           vmResp.Uuid,
					Name:           "test001",
					Vcpus:          1,
					MemoryKib:      1 * 1024 * 1024,
					HypervisorName: hypervisorName,
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
