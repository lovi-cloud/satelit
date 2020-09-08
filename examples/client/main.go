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

	pb "github.com/whywaita/satelit/api/satelit"
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

	fmt.Println("GetImages")
	resp, err := client.ListImage(ctx, &pb.ListImageRequest{})
	if err != nil {
		return err
	}
	for _, i := range resp.Images {
		fmt.Printf("%+v\n", i)
	}

	return nil
}

// SampleStartVirtualMachine is sample of AddVirtualMachine / StartVirtualMachine
func SampleStartVirtualMachine(ctx context.Context, client pb.SatelitClient, imageUUID string) error {
	fmt.Println("AddVirtualMachine")
	resp1, err := client.AddVirtualMachine(ctx, &pb.AddVirtualMachineRequest{
		Name:           "cirros-boot-test",
		Vcpus:          1,
		MemoryKib:      2 * 1024 * 1024,
		RootVolumeGb:   32,
		SourceImageId:  imageUUID,
		HypervisorName: "hv001",
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

	image, err := UploadImage(ctx, client, f, name, "md5:"+hex.EncodeToString(hb))
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
func UploadImage(ctx context.Context, client pb.SatelitClient, src io.Reader, name, description string) (*pb.Image, error) {
	stream, err := client.UploadImage(ctx)
	if err != nil {
		return nil, err
	}

	return uploadImage(stream, src, name, description)
}

func uploadImage(stream pb.Satelit_UploadImageClient, src io.Reader, name, description string) (*pb.Image, error) {
	meta := &pb.UploadImageRequest{
		Value: &pb.UploadImageRequest_Meta{
			Meta: &pb.UploadImageRequestMeta{
				Name:        name,
				Description: description,
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
