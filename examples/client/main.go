package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	pb "github.com/lovi-cloud/satelit/api/satelit"
)

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

	//if err := SampleVolumeOperation(ctx, client); err != nil {
	//	return err
	//}

	//if err := SampleUploadImage(ctx, client); err != nil {
	//	return err
	//}

	//fmt.Println("GetImages")
	//resp, err := client.ListImage(ctx, &pb.ListImageRequest{})
	//if err != nil {
	//	return err
	//}
	//for _, i := range resp.Images {
	//	fmt.Printf("%+v\n", i)
	//}
	//
	if err := SampleStartVirtualMachine(ctx, client, "00000000-0000-0000-0000-000000000000"); err != nil {
		return err
	}

	return nil
}

// SampleVolumeOperation is sample of Volume Add / Delete
func SampleVolumeOperation(ctx context.Context, client pb.SatelitClient) error {
	fmt.Println("AddVolume")
	addVolumeResp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
		Name:             "00000000-0000-0000-0000-000000000000",
		CapacityGigabyte: 20,
		BackendName:      "europa001",
	})
	if err != nil {
		return err
	}
	fmt.Println(addVolumeResp.Volume)

	fmt.Println("GetVolume")
	getVolumeResp, err := client.ListVolume(ctx, &pb.ListVolumeRequest{})
	if err != nil {
		return err
	}
	for _, v := range getVolumeResp.Volumes {
		fmt.Printf("%+v\n", v)
	}

	fmt.Println("DeleteVolume")
	if _, err := client.DeleteVolume(ctx, &pb.DeleteVolumeRequest{
		Id: addVolumeResp.Volume.Id,
	}); err != nil {
		return err
	}

	return nil
}

// SampleStartVirtualMachine is sample of AddVirtualMachine / StartVirtualMachine
func SampleStartVirtualMachine(ctx context.Context, client pb.SatelitClient, imageUUID string) error {
	fmt.Println("AddVirtualMachine")
	resp1, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:              os.Args[1],
		Vcpus:             1,
		MemoryKib:         2 * 1024 * 1024,
		RootVolumeGb:      32,
		SourceImageId:     imageUUID,
		HypervisorName:    "hv0001",
		EuropaBackendName: "europa001",
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

// SampleUploadImage is sample of UploadImage / GetImages / DeleteImage
func SampleUploadImage(ctx context.Context, client pb.SatelitClient) error {
	fmt.Println("UploadImage")
	args := os.Args
	fmt.Printf("args: %s\n", args)
	imageFile := args[1]
	name := filepath.Base(imageFile[:len(imageFile)-len(filepath.Ext(imageFile))])
	f, err := os.Open(imageFile)
	if err != nil {
		return err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	hb := h.Sum(nil)[:16]

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	image, err := UploadImage(ctx, client, f, name, "md5:"+hex.EncodeToString(hb), args[2])
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", image)

	fmt.Println("GetImages")
	resp, err := client.ListImage(ctx, &pb.ListImageRequest{})
	if err != nil {
		return err
	}
	for _, i := range resp.Images {
		fmt.Printf("%+v\n", i)
	}

	//fmt.Println("DeleteImage")
	//deleteResp, err := client.DeleteImage(ctx, &pb.DeleteImageRequest{Id: image.Id})
	//if err != nil {
	//	return err
	//}
	//fmt.Println(deleteResp)

	return nil
}

// UploadImage upload image
func UploadImage(ctx context.Context, client pb.SatelitClient, src io.Reader, name, description, europaBackend string) (*pb.Image, error) {
	stream, err := client.UploadImage(ctx)
	if err != nil {
		return nil, err
	}

	return uploadImage(stream, src, name, description, europaBackend)
}

func uploadImage(stream pb.Satelit_UploadImageClient, src io.Reader, name, description, europaBackend string) (*pb.Image, error) {
	meta := &pb.UploadImageRequest{
		Value: &pb.UploadImageRequest_Meta{
			Meta: &pb.UploadImageRequestMeta{
				Name:              name,
				Description:       description,
				EuropaBackendName: europaBackend,
			}}}
	err := stream.Send(meta)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	for {
		_, err := src.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data := &pb.UploadImageRequest{
			Value: &pb.UploadImageRequest_Chunk{
				Chunk: &pb.UploadImageRequestChunk{
					Data: buf}}}
		err = stream.Send(data)
		if err != nil {
			return nil, err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return resp.Image, nil
}
