package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"sync"

	"github.com/whywaita/satelit/pkg/datastore"

	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/qcow2"
	"github.com/whywaita/satelit/pkg/europa"
)

// A SatelitServer is definition of Satlite API Server
type SatelitServer struct {
	pb.UnimplementedSatelitServer

	Europa    europa.Europa
	Datastore datastore.Datastore
}

// Run start gRPC Server
func (s *SatelitServer) Run() int {
	logger.Logger.Info(fmt.Sprintf("Run satelit server, listen on %s", config.GetValue().API.Listen))
	lis, err := net.Listen("tcp", config.GetValue().API.Listen)
	if err != nil {
		logger.Logger.Error(err.Error())
		return 1
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSatelitServer(grpcServer, s)

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Logger.Error(err.Error())
		return 1
	}

	return 0
}

// GetVolumes call ListVolume to Europa Backend
func (s *SatelitServer) GetVolumes(ctx context.Context, req *pb.GetVolumesRequest) (*pb.GetVolumesResponse, error) {
	volumes, err := s.Europa.ListVolume(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of volume: %w", err)
	}

	var pvs []*pb.Volume
	for _, v := range volumes {
		pvs = append(pvs, v.ToPb())
	}

	return &pb.GetVolumesResponse{
		Volumes: pvs,
	}, nil
}

// AddVolume call CreateVolume to Europa backend
func (s *SatelitServer) AddVolume(ctx context.Context, req *pb.AddVolumeRequest) (*pb.AddVolumeResponse, error) {
	u, err := s.parseRequestUUID(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request id (ID: %s): %w", req.Name, err)
	}

	volume, err := s.Europa.CreateVolumeRaw(ctx, u, int(req.CapacityByte))
	if err != nil {
		return nil, fmt.Errorf("failed to create volume (ID: %s): %w", req.Name, err)
	}

	return &pb.AddVolumeResponse{
		Volume: volume.ToPb(),
	}, nil
}

// AttachVolume call AttachVolume to Europa backend
func (s *SatelitServer) AttachVolume(ctx context.Context, req *pb.AttachVolumeRequest) (*pb.AttachVolumeResponse, error) {
	err := s.Europa.AttachVolume(ctx, req.Id, req.Hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to attach volume to %s (ID: %s): %w", req.Hostname, req.Id, err)
	}

	return &pb.AttachVolumeResponse{}, nil
}

// parseRequestUUID return uuid.UUID from gRPC request string
func (s *SatelitServer) parseRequestUUID(reqName string) (uuid.UUID, error) {
	u := uuid.FromStringOrNil(reqName)
	if u == uuid.Nil {
		return uuid.Nil, fmt.Errorf("failed to parse uuid from string (name: %s)", reqName)
	}

	return u, nil
}

// UploadImage upload to europa backend
func (s *SatelitServer) UploadImage(stream pb.Satelit_UploadImageServer) error {
	logger.Logger.Info("starting UploadImage")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	m, err := s.receiveImage(stream, buf)
	if err != nil {
		return fmt.Errorf("failed to receive image file: %w", err)
	}
	logger.Logger.Debug(fmt.Sprintf("received image (name: %s)", m.name))

	// validate qcow2 image
	b := buf.Bytes()
	reader := bytes.NewReader(b)
	isQcow2, header := qcow2.Probe(reader)
	if isQcow2 == false {
		return errors.New("failed to validate qcow2 image")
	}

	// send to europa
	image, err := s.Europa.UploadImage(ctx, b, m.name, m.description, sanitizeImageSize(header.Size))
	if err != nil {
		return fmt.Errorf("failed to upload image to europa: %w", err)
	}
	logger.Logger.Debug("uploaded image to europa")

	err = stream.SendAndClose(&pb.UploadImageResponse{Image: image.ToPb()})
	if err != nil {
		return fmt.Errorf("failed to send and close: %s", err)
	}
	logger.Logger.Debug("close UploadImage stream")

	// save to image info in database
	err = s.Datastore.PutImage(*image)
	if err != nil {
		return fmt.Errorf("failed to put image to datastore: %w", err)
	}
	logger.Logger.Debug("completed write image to datastore")

	logger.Logger.Info(fmt.Sprintf("UploadImage is successfully! (name: %s)", m.name))
	return nil
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024*64))
	},
}

type meta struct {
	name        string
	description string
}

// header.Size is byte
func sanitizeImageSize(headerSize uint64) int {
	const (
		BYTE float64 = 1 << (10 * iota)
		KILOBYTE
		MEGABYTE
		GIGABYTE
	)

	sizeGB := float64(headerSize) / GIGABYTE
	gb := math.Trunc(sizeGB)

	if sizeGB == gb {
		// not digit loss
		return int(gb)
	}

	// if occurred digit loss, disk capacity need to extend + 1GB
	return int(gb) + 1
}

func (s *SatelitServer) receiveImage(stream pb.Satelit_UploadImageServer, w io.Writer) (meta, error) {
	m := meta{}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return meta{}, fmt.Errorf("failed to recv file: %w", err)
		}

		if mt := resp.GetMeta(); mt != nil {
			m.name = mt.Name
			m.description = mt.Description
		}
		if chunk := resp.GetChunk(); chunk != nil {
			_, err := w.Write(chunk.Data)
			if err != nil {
				return meta{}, fmt.Errorf("failed to write chunk data: %w", err)
			}
		}
	}

	//err := stream.SendAndClose(&pb.UploadImageResponse{})
	//if err != nil {
	//	return meta{}, fmt.Errorf("failed to close stream")
	//}

	return m, nil
}
