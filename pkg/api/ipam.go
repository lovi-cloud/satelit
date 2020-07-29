package api

import (
	"context"

	pb "github.com/whywaita/satelit/api/satelit"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateSubnet create a subnet
func (s *SatelitServer) CreateSubnet(ctx context.Context, req *pb.CreateSubnetRequest) (*pb.CreateSubnetResponse, error) {
	subnet, err := s.IPAM.CreateSubnet(ctx, req.Name, req.Network, req.Start, req.End, req.Gateway, req.DnsServer, req.MetadataServer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create subnet: %+v", err)
	}

	return &pb.CreateSubnetResponse{
		Subnet: &pb.Subnet{
			Uuid:           subnet.UUID.String(),
			Name:           subnet.Name,
			Network:        subnet.Network.String(),
			Start:          subnet.Start.String(),
			End:            subnet.End.String(),
			Gateway:        subnet.Gateway.String(),
			DnsServer:      subnet.DNSServer.String(),
			MetadataServer: subnet.MetadataServer.String(),
		},
	}, nil
}

// GetSubnet retrieves address according to the parameters given
func (s *SatelitServer) GetSubnet(ctx context.Context, req *pb.GetSubnetRequest) (*pb.GetSubnetResponse, error) {
	u, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request id (need uuid): %+v", err)
	}
	subnet, err := s.IPAM.GetSubnet(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieves subnet: %+v", err)
	}

	return &pb.GetSubnetResponse{
		Subnet: &pb.Subnet{
			Uuid:           subnet.UUID.String(),
			Name:           subnet.Name,
			Network:        subnet.Network.String(),
			Start:          subnet.Start.String(),
			End:            subnet.End.String(),
			Gateway:        subnet.Gateway.String(),
			DnsServer:      subnet.DNSServer.String(),
			MetadataServer: subnet.MetadataServer.String(),
		},
	}, nil
}

// ListSubnet retrieves all subnets
func (s *SatelitServer) ListSubnet(ctx context.Context, req *pb.ListSubnetRequest) (*pb.ListSubnetResponse, error) {
	subnets, err := s.IPAM.ListSubnet(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list subnet: %+v", err)
	}

	tmp := make([]*pb.Subnet, len(subnets))
	for i, subnet := range subnets {
		tmp[i] = &pb.Subnet{
			Uuid:           subnet.UUID.String(),
			Name:           subnet.Name,
			Network:        subnet.Network.String(),
			Start:          subnet.Start.String(),
			End:            subnet.End.String(),
			Gateway:        subnet.Gateway.String(),
			DnsServer:      subnet.DNSServer.String(),
			MetadataServer: subnet.MetadataServer.String(),
		}
	}

	return &pb.ListSubnetResponse{
		Subnets: tmp,
	}, nil
}

// DeleteSubnet deletes a subnet
func (s *SatelitServer) DeleteSubnet(ctx context.Context, req *pb.DeleteSubnetRequest) (*pb.DeleteSubnetResponse, error) {
	u, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request id (need uuid): %+v", err)
	}
	if err := s.IPAM.DeleteSubnet(ctx, u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete subnet: %+v", err)
	}

	return &pb.DeleteSubnetResponse{}, nil
}

// CreateAddress create a address
func (s *SatelitServer) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	u, err := s.parseRequestUUID(req.SubnetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request subnet id (need uuid): %+v", err)
	}
	address, err := s.IPAM.CreateAddress(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create address: %+v", err)
	}

	return &pb.CreateAddressResponse{
		Address: &pb.Address{
			Uuid:     address.UUID.String(),
			Ip:       address.IP.String(),
			SubnetId: address.SubnetID.String(),
		},
	}, nil
}

// GetAddress retrieves address according to the parameters given
func (s *SatelitServer) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	u, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request uuid (need uuid): %+v", err)
	}
	address, err := s.IPAM.GetAddress(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get address: %+v", err)
	}

	return &pb.GetAddressResponse{
		Address: &pb.Address{
			Uuid:     address.UUID.String(),
			Ip:       address.IP.String(),
			SubnetId: address.SubnetID.String(),
		},
	}, nil
}

// ListAddress retrieves all address according to the parameters given.
func (s *SatelitServer) ListAddress(ctx context.Context, req *pb.ListAddressRequest) (*pb.ListAddressResponse, error) {
	u, err := s.parseRequestUUID(req.SubnetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request subnet id (need uuid): %+v", err)
	}
	addresses, err := s.IPAM.ListAddressBySubnetID(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list address: %+v", err)
	}

	tmp := make([]*pb.Address, len(addresses))
	for i, address := range addresses {
		tmp[i] = &pb.Address{
			Uuid:     address.UUID.String(),
			Ip:       address.IP.String(),
			SubnetId: address.SubnetID.String(),
		}
	}

	return &pb.ListAddressResponse{
		Addresses: tmp,
	}, nil
}

// DeleteAddress deletes address
func (s *SatelitServer) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	u, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request uuid (need uuid): %+v", err)
	}
	if err := s.IPAM.DeleteAddress(ctx, u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete address: %+v", err)
	}

	return &pb.DeleteAddressResponse{}, nil
}

// CreateLease create a lease.
func (s *SatelitServer) CreateLease(ctx context.Context, req *pb.CreateLeaseRequest) (*pb.CreateLeaseResponse, error) {
	u, err := s.parseRequestUUID(req.AddressId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request address id (need uuid): %+v", err)
	}
	lease, err := s.IPAM.CreateLease(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create lease: %+v", err)
	}

	return &pb.CreateLeaseResponse{
		Lease: &pb.Lease{
			MacAddress: lease.MacAddress.String(),
			AddressId:  lease.AddressID.String(),
		},
	}, nil
}

// GetLease retrieves address according to the parameters given.
func (s *SatelitServer) GetLease(ctx context.Context, req *pb.GetLeaseRequest) (*pb.GetLeaseResponse, error) {
	leaseID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request lease id: %+v", err)
	}
	lease, err := s.IPAM.GetLease(ctx, leaseID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to retrieve lease: %+v", err)
	}

	return &pb.GetLeaseResponse{
		Lease: &pb.Lease{
			MacAddress: lease.MacAddress.String(),
			AddressId:  lease.AddressID.String(),
		},
	}, nil
}

// ListLease retrieves all leases according to the parameters given.
func (s *SatelitServer) ListLease(ctx context.Context, req *pb.ListLeaseRequest) (*pb.ListLeaseResponse, error) {
	leases, err := s.IPAM.ListLease(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list leases: %+v", err)
	}

	tmp := make([]*pb.Lease, len(leases))
	for i, lease := range leases {
		tmp[i] = &pb.Lease{
			MacAddress: lease.MacAddress.String(),
			AddressId:  lease.AddressID.String(),
		}
	}

	return &pb.ListLeaseResponse{
		Leases: tmp,
	}, nil
}

// DeleteLease deletes lease
func (s *SatelitServer) DeleteLease(ctx context.Context, req *pb.DeleteLeaseRequest) (*pb.DeleteLeaseResponse, error) {
	leaseID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, err
	}
	if err := s.IPAM.DeleteLease(ctx, leaseID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete lease: %+v", err)
	}

	return &pb.DeleteLeaseResponse{}, nil
}
