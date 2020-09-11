package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"github.com/whywaita/satelit/pkg/europa"

	"github.com/whywaita/satelit/pkg/scheduler/scheduler"

	"github.com/whywaita/satelit/pkg/datastore"

	dspb "github.com/whywaita/satelit/api/satelit_datastore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/testutils"
	datastoreMemory "github.com/whywaita/satelit/pkg/datastore/memory"
	europaMemory "github.com/whywaita/satelit/pkg/europa/memory"
	"github.com/whywaita/satelit/pkg/ganymede"
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
func NewMemorySatelit(ds datastore.Datastore) *SatelitServer {
	ipamBackend := ipam.New(ds)
	e := europaMemory.New(ds)
	ganymede := ganymedeMemory.New(ds)
	sche := scheduler.New(ds)

	return &SatelitServer{
		Europa: map[string]europa.Europa{
			"memory": e,
		},
		IPAM:      ipamBackend,
		Datastore: ds,
		Ganymede:  ganymede,
		Scheduler: sche,
	}
}

// NewMemorySatelitDatastore create in-memory Satelit datastore Server
// for testing Satelit Datastore API
func NewMemorySatelitDatastore(ds datastore.Datastore) *SatelitDatastore {
	return &SatelitDatastore{
		Datastore: ds,
	}
}

const bufSize = 1024 * 1024

var lisSatelit *bufconn.Listener
var resetAll func()
var lisDatastore *bufconn.Listener

func init() {
	logger.New("debug")

	ds := datastoreMemory.New()

	server := NewMemorySatelit(ds)

	lisSatelit = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSatelitServer(s, server)
	go func() {
		if err := s.Serve(lisSatelit); err != nil {
			log.Fatal(err)
		}
	}()

	dsServer := NewMemorySatelitDatastore(ds)
	lisDatastore = bufconn.Listen(bufSize)
	sDs := grpc.NewServer()
	dspb.RegisterSatelitDatastoreServer(sDs, dsServer)
	go func() {
		if err := sDs.Serve(lisDatastore); err != nil {
			log.Fatal(err)
		}
	}()

	resetAll = func() {
		ds := datastoreMemory.New()

		server.Datastore = ds
		server.IPAM = ipam.New(ds)
		server.Europa = map[string]europa.Europa{
			"memory": europaMemory.New(ds),
		}
		server.Ganymede = ganymedeMemory.New(ds)
		server.Scheduler = scheduler.New(ds)

		dsServer.Datastore = ds
	}
}
func bufDialerSatelit(ctx context.Context, address string) (net.Conn, error) {
	return lisSatelit.Dial()
}

func bufDialerDatastore(ctx context.Context, address string) (net.Conn, error) {
	return lisDatastore.Dial()
}

func getSatelitClient() (context.Context, pb.SatelitClient, func() error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialerSatelit), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewSatelitClient(conn)

	return ctx, client, func() error {
		resetAll()
		return conn.Close()
	}
}

func getDatastoreClient() (context.Context, dspb.SatelitDatastoreClient, func() error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialerDatastore), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := dspb.NewSatelitDatastoreClient(conn)

	return ctx, client, func() error {
		resetAll()
		return conn.Close()
	}
}

func setupTeleskop(nodes []ganymede.NUMANode) (string, func(), error) {
	ctx, client, _ := getDatastoreClient()
	hypervisorName := "dummy"
	iqn := "dummy-iqn"

	var ep string
	ep, teardown, err := testutils.NewDummyTeleskop()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create dummy teleskop: %w", err)
	}
	err = teleskop.New(map[string]string{hypervisorName: ep})
	if err != nil {
		return "", nil, fmt.Errorf("failed to teleskop.New: %w", err)
	}

	if nodes != nil {
		var pbNodes []*dspb.NumaNode
		for _, n := range nodes {
			pbNodes = append(pbNodes, n.ToPb())
		}

		if _, err := client.RegisterTeleskopAgent(ctx, &dspb.RegisterTeleskopAgentRequest{
			Hostname: hypervisorName,
			Endpoint: ep,
			Iqn:      iqn,
			Nodes:    pbNodes,
		}); err != nil {
			return "", nil, fmt.Errorf("failed to RegisterTeleskopAgent: %w", err)
		}
	}

	t := func() {
		teardown()
		//teardownDS()
	}

	return hypervisorName, t, nil
}

func uploadImage(ctx context.Context, client pb.SatelitClient, image io.Reader) (*pb.UploadImageResponse, error) {
	stream, err := client.UploadImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to call upload image: %w", err)
	}

	err = stream.Send(&pb.UploadImageRequest{
		Value: &pb.UploadImageRequest_Meta{
			Meta: &pb.UploadImageRequestMeta{
				Name:              "image001",
				Description:       "desc",
				EuropaBackendName: "memory",
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
