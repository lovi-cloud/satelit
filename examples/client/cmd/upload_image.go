package cmd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"github.com/spf13/cobra"

	pb "github.com/lovi-cloud/satelit/api/satelit"
)

var (
	pathImage         string
	europaBackendName string
)

var uploadImageCmd = &cobra.Command{
	Use:   "upload-image",
	Short: "Upload image to satelit from qcow2 file",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(satelitAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			return err
		}

		client := pb.NewSatelitClient(conn)
		return SampleUploadImage(context.Background(), client)
	},
}

func init() {
	uploadImageCmd.Flags().StringVarP(&pathImage, "image", "", "", "file path of qcow2 image")
	uploadImageCmd.MarkFlagRequired("image")
	uploadImageCmd.Flags().StringVarP(&europaBackendName, "backend", "", "", "backend name of europa")
	uploadImageCmd.MarkFlagRequired("backend")

	rootCmd.AddCommand(uploadImageCmd)
}

// SampleUploadImage is sample of UploadImage / GetImages / DeleteImage
func SampleUploadImage(ctx context.Context, client pb.SatelitClient) error {
	fmt.Println("UploadImage")
	name := filepath.Base(pathImage[:len(pathImage)-len(filepath.Ext(pathImage))])
	f, err := os.Open(pathImage)
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

	image, err := UploadImage(ctx, client, f, name, "md5:"+hex.EncodeToString(hb), europaBackendName)
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
