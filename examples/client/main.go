package main

import (
	"fmt"
	"os"

	"github.com/lovi-cloud/satelit/examples/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

//// SampleVolumeOperation is sample of Volume Add / Delete
//func SampleVolumeOperation(ctx context.Context, client pb.SatelitClient) error {
//	fmt.Println("AddVolume")
//	addVolumeResp, err := client.AddVolume(ctx, &pb.AddVolumeRequest{
//		Name:             "00000000-0000-0000-0000-000000000000",
//		CapacityGigabyte: 20,
//		BackendName:      "europa001",
//	})
//	if err != nil {
//		return err
//	}
//	fmt.Println(addVolumeResp.Volume)
//
//	fmt.Println("GetVolume")
//	getVolumeResp, err := client.ListVolume(ctx, &pb.ListVolumeRequest{})
//	if err != nil {
//		return err
//	}
//	for _, v := range getVolumeResp.Volumes {
//		fmt.Printf("%+v\n", v)
//	}
//
//	fmt.Println("DeleteVolume")
//	if _, err := client.DeleteVolume(ctx, &pb.DeleteVolumeRequest{
//		Id: addVolumeResp.Volume.Id,
//	}); err != nil {
//		return err
//	}
//
//	return nil
//}
