package api

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/whywaita/satelit/api/satelit_datastore"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
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
	grpcServer := grpc.NewServer()
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

	lease, err := s.Datastore.GetDHCPLeaseByMACAddress(ctx, mac)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get DHCP lease information: %+v", err)
	}

	return &pb.GetDHCPLeaseResponse{
		Lease: &pb.DHCPLease{
			MacAddress:     lease.MacAddress.String(),
			Ip:             lease.IP.String(),
			Gateway:        lease.Gateway.String(),
			DnsServer:      lease.DNSServer.String(),
			MetadataServer: lease.MetadataServer.String(),
		},
	}, nil
}
