package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/satelit/internal/client/teleskop"

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
	"github.com/whywaita/satelit/pkg/ganymede"
)

// A SatelitDatastore is definition of Satelit Datastore API Server
type SatelitDatastore struct {
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

// GetHostnameByAddress is
func (s *SatelitDatastore) GetHostnameByAddress(ctx context.Context, req *pb.GetHostnameByAddressRequest) (*pb.GetHostnameByAddressResponse, error) {
	address := net.ParseIP(req.Address)
	if address == nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request address")
	}

	hostname, err := s.Datastore.GetHostnameByAddress(types.IP(address))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get hostname by address: %+v", err)
	}

	return &pb.GetHostnameByAddressResponse{
		Hostname: hostname,
	}, nil
}

// GetISUCONUserKeys is
func (s *SatelitDatastore) GetISUCONUserKeys(ctx context.Context, req *pb.GetISUCONUserKeysRequest) (*pb.GetISUCONUserKeysResponse, error) {
	address := net.ParseIP(req.Address)
	if address == nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request address")
	}

	return &pb.GetISUCONUserKeysResponse{
		Keys: strings.Split(strings.TrimSpace(adminKeys), "\n"),
	}, nil
}

// GetISUCONAdminKeys is
func (s *SatelitDatastore) GetISUCONAdminKeys(ctx context.Context, req *pb.GetISUCONAdminKeysRequest) (*pb.GetISUCONAdminKeysResponse, error) {
	return &pb.GetISUCONAdminKeysResponse{
		Keys: strings.Split(strings.TrimSpace(adminKeys), "\n"),
	}, nil
}

// ListBridge is
func (s *SatelitDatastore) ListBridge(ctx context.Context, req *pb.ListBridgeRequest) (*pb.ListBridgeResponse, error) {
	bridges, err := s.Datastore.ListBridge(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bridge list: %+v", err)
	}

	resp := make([]*pb.ListBridgeResponse_Bridge, len(bridges))
	for i, bridge := range bridges {
		resp[i] = &pb.ListBridgeResponse_Bridge{
			Name:            bridge.Name,
			VlanId:          bridge.VLANID,
			ParentInterface: "bond0",
		}
		if bridge.VLANID == 0 {
			resp[i].MetadataCidr = ""
			resp[i].InternalOnly = true
			continue
		}

		subnet, err := s.Datastore.GetSubnetByVLAN(ctx, bridge.VLANID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get subnet by vlan=%d: %+v", bridge.VLANID, err)
		}
		mask, _ := subnet.Network.Mask.Size()
		resp[i].MetadataCidr = fmt.Sprintf("%s/%d", subnet.MetadataServer.String(), mask)
		resp[i].InternalOnly = false
	}

	return &pb.ListBridgeResponse{
		Bridges: resp,
	}, nil
}

// RegisterTeleskopAgent register new teleskop agent
func (s *SatelitDatastore) RegisterTeleskopAgent(ctx context.Context, req *pb.RegisterTeleskopAgentRequest) (*pb.RegisterTeleskopAgentResponse, error) {
	err := teleskop.AddClient(req.Hostname, req.Endpoint)
	if err != nil && !errors.Is(err, teleskop.ErrTeleskopAgentAlreadyExist) {
		return nil, status.Errorf(codes.InvalidArgument, "failed to register teleskop agent: %+v", err)
	}

	var nodes []ganymede.NUMANode
	for _, n := range req.Nodes {
		var pairs []ganymede.CorePair
		for _, p := range n.Pairs {
			pair := ganymede.CorePair{
				UUID:         uuid.NewV4(),
				PhysicalCore: p.PhysicalCore,
				LogicalCore:  p.LogicalCore,
			}
			pairs = append(pairs, pair)
		}

		node := ganymede.NUMANode{
			UUID:            uuid.NewV4(),
			CorePairs:       pairs,
			PhysicalCoreMin: n.PhysicalCoreMin,
			PhysicalCoreMax: n.PhysicalCoreMax,
			LogicalCoreMin:  n.LogicalCoreMin,
			LogicalCoreMax:  n.LogicalCoreMax,
		}
		nodes = append(nodes, node)
	}

	hypervisorID, err := s.Datastore.PutHypervisor(ctx, req.Iqn, req.Hostname)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write hypervisor to datastore: %+v", err)
	}

	if err := s.Datastore.PutHypervisorNUMANode(ctx, nodes, hypervisorID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write hypervisor cores to datastore: %+v", err)
	}

	return &pb.RegisterTeleskopAgentResponse{}, nil
}

// GetCPUCoreByPinningGroup retrieve pinned cpu pair.
func (s *SatelitDatastore) GetCPUCoreByPinningGroup(ctx context.Context, req *pb.GetCPUCoreByPinningGroupRequest) (*pb.GetCPUCoreByPinningGroupResponse, error) {
	cpg, err := s.Datastore.GetCPUPinningGroupByName(ctx, req.PinningGroupName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cpu pinning group: %+v", err)
	}

	pinneds, err := s.Datastore.GetPinnedCoreByPinningGroup(ctx, cpg.UUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cpu pinned: %+v", err)
	}

	var pairs []*pb.CorePair
	for _, pinned := range pinneds {
		pair, err := s.Datastore.GetCPUCorePair(ctx, pinned.CorePairID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get pinned cpu corepairs: %+v", err)
		}
		pairs = append(pairs, pair.ToPb())
	}

	return &pb.GetCPUCoreByPinningGroupResponse{
		Pairs: pairs,
	}, nil
}
