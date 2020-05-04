package api

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
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

// GetVolumes response volumes to pb
func (s *SatelitServer) GetVolumes(ctx context.Context, req *pb.GetVolumesRequest) (*pb.GetVolumesResponse, error) {
	volumes, err := s.Europa.ListVolume(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of volume")
	}

	var pvs []*pb.Volume
	for _, v := range volumes {
		pvs = append(pvs, v.ToPb())
	}

	return &pb.GetVolumesResponse{
		Volumes: pvs,
	}, nil
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
