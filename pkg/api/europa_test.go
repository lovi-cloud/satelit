package api

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
	"testing"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/testutils"
)

const gzipCompressedQcow2String = `
H4sIAKNA+14AA+3OPQ6CQBAG0F3wAB6B05hYWlkraiAxYlYo8BQe18RG/Cvs6Gjem8xMppjkWy2W
jxBCHv7NPyv73fE9878rhnFmQ1dVfl8P+xbCrk5tX2zrdtx3LJuUunM78iOG4+baF2l/KJvu1F5G
ZgQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAOAlmzoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAk4vfAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgMETEqcE5ggAAwA=
`

func TestSatelitServer_AddVirtualMachine(t *testing.T) {
	ep, teardownTeleskop, err := testutils.NewDummyTeleskop()
	if err != nil {
		t.Fatalf("failed to start dummy teleskop: %+v\n", err)
	}
	defer teardownTeleskop()

	if err := teleskop.New(map[string]string{"dummy": ep}); err != nil {
		t.Fatalf("failed to create teleskop client: %+v\n", err)
	}

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	stream, err := client.UploadImage(ctx)
	if err != nil {
		t.Fatalf("failed to upload image: %+v\n", err)
	}

	err = stream.Send(&pb.UploadImageRequest{
		Value: &pb.UploadImageRequest_Meta{
			Meta: &pb.UploadImageRequestMeta{
				Name:        "image001",
				Description: "desc",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to send meta data: %+v\n", err)
	}

	dummyImage, err := getDummyQcow2Image()
	if err != nil {
		t.Fatalf("failed to get dummy qcow2 image: %+v\n", err)
	}

	buff := make([]byte, 1024)
	for {
		n, err := dummyImage.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			t.Fatalf("failed to read dummy image: %+v\n", err)
		}
		err = stream.Send(&pb.UploadImageRequest{
			Value: &pb.UploadImageRequest_Chunk{
				Chunk: &pb.UploadImageRequestChunk{
					Data: buff[:n],
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to send data: %+v\n", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("failed to close and recv stream: %+v\n", err)
	}

	_, err = client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:           "test001",
		Vcpus:          1,
		MemoryKib:      1 * 1024 * 1024,
		RootVolumeGb:   10,
		SourceImageId:  resp.Image.Id,
		HypervisorName: "dummy",
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
