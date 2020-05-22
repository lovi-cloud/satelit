package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"

	"github.com/whywaita/satelit/internal/config"

	pb "github.com/whywaita/satelit/api/satelit"
)

func init() {
	conf := "./../../configs/satelit.yaml"
	err := config.Load(&conf)
	if err != nil {
		panic(err)
	}
}

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
	resp, err := client.GetImages(ctx, &pb.GetImagesRequest{})
	if err != nil {
		return err
	}
	for _, i := range resp.Images {
		fmt.Printf("%+v\n", i)
	}

	fmt.Println("UploadImage")
	args := os.Args
	fmt.Printf("args: %s\n", args)
	imageFile := args[1]
	f, err := os.Open(imageFile)
	if err != nil {
		return err
	}

	image, err := UploadImage(ctx, client, f, "whywaita-test-cirros", "cirros")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", image)

	fmt.Println("GetImages")
	resp, err = client.GetImages(ctx, &pb.GetImagesRequest{})
	if err != nil {
		return err
	}
	for _, i := range resp.Images {
		fmt.Printf("%+v\n", i)
	}

	fmt.Println("DeleteImage")
	deleteResp, err := client.DeleteImage(ctx, &pb.DeleteImageRequest{Id: image.Id})
	if err != nil {
		return err
	}
	fmt.Println(deleteResp)

	//fmt.Println("AddVolume")
	//resp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
	//	Name:         u.String(),
	//	CapacityByte: 23,
	//})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%+v\n", resp.Volume)
	//
	//var hostname string
	//for h := range config.GetValue().Teleskop.Endpoints {
	//	hostname = h // set last hostname
	//}
	//
	//fmt.Println("AttachVolume")
	//_, err = client.AttachVolume(ctx, &pb.AttachVolumeRequest{
	//	Id:       resp.Volume.Id,
	//	Hostname: hostname,
	//})
	//if err != nil {
	//	return err
	//}

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
