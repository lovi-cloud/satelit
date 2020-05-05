package api

import (
	"context"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/europa"
)

// A SatelitServer is definition of Satlite API Server
type SatelitServer struct {
	pb.UnimplementedSatelitServer

	Europa europa.Europa
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

	volume, err := s.Europa.CreateVolume(ctx, u, int(req.CapacityByte))
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
