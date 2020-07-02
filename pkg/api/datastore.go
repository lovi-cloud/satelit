package api

import (
	"context"
	"fmt"
	"net"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/whywaita/satelit/api/satelit_datastore"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/mysql/types"
	"github.com/whywaita/satelit/pkg/datastore"
)

// A SatelitDatastore is definition of Satelit Datastore API Server
type SatelitDatastore struct {
	pb.UnimplementedSatelitDatastoreServer

	Datastore datastore.Datastore
}

// Run start gRPC Server
func (s *SatelitDatastore) Run() error {
	logger.Logger.Info(fmt.Sprintf("Run satelit server, listen on %s", config.GetValue().Datastore.Listen))
	lis, err := net.Listen("tcp", config.GetValue().Datastore.Listen)
	if err != nil {
		return err
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger.Logger)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			}),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.PayloadStreamServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				if strings.Contains(fullMethodName, "UploadImage") {
					return false
				}
				return true
			}),
		),
	)
	pb.RegisterSatelitDatastoreServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

// GetDHCPLease return DHCP lease information.
func (s *SatelitDatastore) GetDHCPLease(ctx context.Context, req *pb.GetDHCPLeaseRequest) (*pb.GetDHCPLeaseResponse, error) {
	mac, err := net.ParseMAC(req.MacAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse MAC address: %+v", err)
	}

	lease, err := s.Datastore.GetDHCPLeaseByMACAddress(ctx, types.HardwareAddr(mac))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get DHCP lease information: %+v", err)
	}

	return &pb.GetDHCPLeaseResponse{
		Lease: &pb.DHCPLease{
			MacAddress:     lease.MacAddress.String(),
			Ip:             lease.IP.String(),
			Network:        lease.Network.String(),
			Gateway:        lease.Gateway.String(),
			DnsServer:      lease.DNSServer.String(),
			MetadataServer: lease.MetadataServer.String(),
		},
	}, nil
}
