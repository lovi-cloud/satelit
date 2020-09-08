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

func TestSatelitServer_ShowVolume(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: 1,
	})
	if err != nil {
		t.Fatalf("failed to create test volume: %+v\n", err)
	}

	tests := []struct {
		input *pb.ShowVolumeRequest
		want  *pb.ShowVolumeResponse
		err   bool
	}{
		{
			input: &pb.ShowVolumeRequest{
				Id: resp.Volume.Id,
			},
			want: &pb.ShowVolumeResponse{
				Volume: resp.Volume,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ShowVolume(ctx, test.input)
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

func TestSatelitServer_ListVolume(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: 1,
	})
	if err != nil {
		t.Fatalf("failed to create test volume: %+v\n", err)
	}

	tests := []struct {
		input *pb.ListVolumeRequest
		want  *pb.ListVolumeResponse
		err   bool
	}{
		{
			input: &pb.ListVolumeRequest{},
			want: &pb.ListVolumeResponse{
				Volumes: []*pb.Volume{
					resp.Volume,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListVolume(ctx, test.input)
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

func TestSatelitServer_AddVolumeImage(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	tests := []struct {
		input *pb.AddVolumeImageRequest
		want  *pb.AddVolumeImageResponse
		err   bool
	}{
		{
			input: &pb.AddVolumeImageRequest{
				Name:             testUUID,
				CapacityGigabyte: 10,
				SourceImageId:    imageResp.Image.Id,
			},
			want: &pb.AddVolumeImageResponse{
				Volume: &pb.Volume{
					Id:               testUUID,
					CapacityGigabyte: 10,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.AddVolumeImage(ctx, test.input)
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

func TestSatelitServer_AttachVolume(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	resp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: 10,
	})
	if err != nil {
		t.Fatalf("failed to add test volume: %+v", err)
	}

	tests := []struct {
		input *pb.AttachVolumeRequest
		want  *pb.AttachVolumeResponse
		err   bool
	}{
		{
			input: &pb.AttachVolumeRequest{
				Id:       resp.Volume.Id,
				Hostname: hypervisorName,
			},
			want: &pb.AttachVolumeResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.AttachVolume(ctx, test.input)
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

func TestSatelitServer_DetachVolume(t *testing.T) {
	hypervisorName, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	volume, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: 10,
	})
	if err != nil {
		t.Fatalf("failed to add test volume: %+v\n", err)
	}

	_, err = client.AttachVolume(ctx, &pb.AttachVolumeRequest{
		Id:       volume.Volume.Id,
		Hostname: hypervisorName,
	})
	if err != nil {
		t.Fatalf("failed to attach test volume: %+v\n", err)
	}

	tests := []struct {
		input *pb.DetachVolumeRequest
		want  *pb.DetachVolumeResponse
		err   bool
	}{
		{
			input: &pb.DetachVolumeRequest{
				Id: volume.Volume.Id,
			},
			want: &pb.DetachVolumeResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DetachVolume(ctx, test.input)
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

func TestSatelitServer_DeleteVolume(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	volume, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             testUUID,
		CapacityGigabyte: 10,
	})
	if err != nil {
		t.Fatalf("failed to add test volume: %+v\n", err)
	}

	tests := []struct {
		input *pb.DeleteVolumeRequest
		want  *pb.DeleteVolumeResponse
		err   bool
	}{
		{
			input: &pb.DeleteVolumeRequest{
				Id: volume.Volume.Id,
			},
			want: &pb.DeleteVolumeResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteVolume(ctx, test.input)
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

func TestSatelitServer_ListImage(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	tests := []struct {
		input *pb.ListImageRequest
		want  *pb.ListImageResponse
		err   bool
	}{
		{
			input: &pb.ListImageRequest{},
			want: &pb.ListImageResponse{
				Images: []*pb.Image{
					imageResp.Image,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListImage(ctx, test.input)
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

func TestSatelitServer_UploadImage(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	image, err := getDummyQcow2Image()
	if err != nil {
		t.Fatalf("failed to get dummy image: %+v\n", err)
	}

	tests := []struct {
		input *pb.UploadImageRequest
		image io.Reader
		want  *pb.UploadImageResponse
		err   bool
	}{
		{
			input: &pb.UploadImageRequest{
				Value: &pb.UploadImageRequest_Meta{
					Meta: &pb.UploadImageRequestMeta{
						Name:        "image001",
						Description: "desc",
					},
				},
			},
			image: image,
			want: &pb.UploadImageResponse{
				Image: &pb.Image{
					Id:          "",
					Name:        "image001",
					VolumeId:    "",
					Description: "desc",
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := uploadImage(ctx, client, test.image)
		if got != nil {
			test.want.Image.Id = got.Image.Id
			test.want.Image.VolumeId = got.Image.VolumeId
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

func TestSatelitServer_DeleteImage(t *testing.T) {
	ctx, client, teardown := getSatelitClient()
	defer teardown()

	imageResp, err := uploadDummyImage(ctx, client)
	if err != nil {
		t.Fatalf("failed to upload dummy image: %+v\n", err)
	}

	tests := []struct {
		input *pb.DeleteImageRequest
		want  *pb.DeleteImageResponse
		err   bool
	}{
		{
			input: &pb.DeleteImageRequest{
				Id: imageResp.Image.Id,
			},
			want: &pb.DeleteImageResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteImage(ctx, test.input)
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
