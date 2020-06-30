package api

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/go-test/deep"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/logger"
	datastoreMemory "github.com/whywaita/satelit/pkg/datastore/memory"
	europaMemory "github.com/whywaita/satelit/pkg/europa/memory"
	ganymedeMemory "github.com/whywaita/satelit/pkg/ganymede/memory"
	"github.com/whywaita/satelit/pkg/ipam/ipam"
)

func TestSanitizeImageSize(t *testing.T) {
	inputs := []uint64{
		536870912,  // 0.5 GB
		1073741824, // 1   GB
		1610612736, // 1.5 GB
	}
	outputs := []int{
		1,
		1,
		2,
	}

	for i, input := range inputs {
		o := sanitizeImageSize(input)
		if outputs[i] != o {
			t.Errorf("failed to sanitize (input: %d)", input)
		}
	}
}

// NewMemorySatelit create in-memory Satelit API Server
// for testing Satelit API
func NewMemorySatelit() *SatelitServer {
	europa := europaMemory.New()
	ds := datastoreMemory.New()
	ipamBackend := ipam.New(ds)
	ganymede := ganymedeMemory.New(ds)

	return &SatelitServer{
		Europa:    europa,
		IPAM:      ipamBackend,
		Datastore: ds,
		Ganymede:  ganymede,
	}
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	logger.New("debug")
	server := NewMemorySatelit()

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSatelitServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}
func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func getSatelitClient() (context.Context, pb.SatelitClient, func() error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewSatelitClient(conn)

	return ctx, client, conn.Close
}

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
