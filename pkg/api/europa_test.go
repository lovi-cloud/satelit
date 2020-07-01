package api

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
	"testing"

	"github.com/go-test/deep"

	pb "github.com/whywaita/satelit/api/satelit"
)

const gzipCompressedQcow2String = `
H4sIAKNA+14AA+3OPQ6CQBAG0F3wAB6B05hYWlkraiAxYlYo8BQe18RG/Cvs6Gjem8xMppjkWy2W
jxBCHv7NPyv73fE9878rhnFmQ1dVfl8P+xbCrk5tX2zrdtx3LJuUunM78iOG4+baF2l/KJvu1F5G
ZgQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAOAlmzoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAk4vfAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgMETEqcE5ggAAwA=
`

const (
	testVolumeName       = "TEST_VOLUME"
	testCapacityGigabyte = 8
	testUUID             = "90dd6cd4-b3e4-47f3-9af5-47f78efc8fc7"
)

func TestSatelitServer_AddVolume(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	req := &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: testCapacityGigabyte,
	}

	resp, err := client.AddVolume(ctx, req)
	if err != nil {
		t.Errorf("AddVolume return error: %+v", err)
	}

	want := pb.Volume{
		Id:               testUUID,
		CapacityGigabyte: testCapacityGigabyte,
	}

	if diff := deep.Equal(resp.Volume, &want); diff != nil {
		t.Error(diff)
	}
}

func TestSatelitServer_AddVirtualMachine(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop()
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	image, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	_, err = client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:           "test001",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   10,
		SourceImageId:  image.Id,
		HypervisorName: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to add virtual machine: %+v\n", err)
	}
}

func getDummyQcow2Image() (io.Reader, error) {
	return gzip.NewReader(
		base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(
				strings.TrimSpace(gzipCompressedQcow2String),
			),
		),
	)
}
