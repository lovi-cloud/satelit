package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/testutils"
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
	ds := datastoreMemory.New()
	ipamBackend := ipam.New(ds)
	europa := europaMemory.New(ds)
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
var resetSatelit func()

func init() {
	logger.New("debug")
	server := NewMemorySatelit()
	resetSatelit = func() {
		server.Datastore = datastoreMemory.New()
		server.IPAM = ipam.New(server.Datastore)
		server.Europa = europaMemory.New(server.Datastore)
		server.Ganymede = ganymedeMemory.New(server.Datastore)
	}

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

	return ctx, client, func() error {
		resetSatelit()
		return conn.Close()
	}
}

func setupTeleskop() (hypervisorName string, teardown func(), err error) {
	hypervisorName = "dummy"

	var ep string
	ep, teardown, err = testutils.NewDummyTeleskop()
	if err != nil {
		return
	}
	err = teleskop.New(map[string]string{hypervisorName: ep})
	if err != nil {
		return
	}
	return hypervisorName, teardown, nil
}

func uploadImage(ctx context.Context, client pb.SatelitClient, image io.Reader) (*pb.UploadImageResponse, error) {
	stream, err := client.UploadImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to call upload image: %w", err)
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
		return nil, fmt.Errorf("failed to send meta data: %w", err)
	}

	buff := make([]byte, 1024)
	for {
		n, err := image.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read image: %w", err)
		}
		err = stream.Send(&pb.UploadImageRequest{
			Value: &pb.UploadImageRequest_Chunk{
				Chunk: &pb.UploadImageRequestChunk{
					Data: buff[:n],
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send data: %w", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to close and recv stream: %w", err)
	}

	return resp, nil
}

func uploadDummyImage(ctx context.Context, client pb.SatelitClient) (*pb.UploadImageResponse, error) {
	dummyImage, err := getDummyQcow2Image()
	if err != nil {
		return nil, fmt.Errorf("failed to get dummy qcow2 image: %w", err)
	}
	return uploadImage(ctx, client, dummyImage)
}
