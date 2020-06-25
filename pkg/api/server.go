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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"

	agentpb "github.com/whywaita/satelit/api"
	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/qcow2"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A SatelitServer is definition of Satlite API Server
type SatelitServer struct {
	pb.SatelitServer

	Europa europa.Europa
	IPAM   ipam.IPAM

	Datastore datastore.Datastore
	Ganymede  ganymede.Ganymede
}

// Run start gRPC Server
func (s *SatelitServer) Run() error {
	logger.Logger.Info(fmt.Sprintf("Run satelit server, listen on %s", config.GetValue().API.Listen))
	lis, err := net.Listen("tcp", config.GetValue().API.Listen)
	if err != nil {
		return err
	}
	opts := []grpc_zap.Option{
		grpc_zap.WithMessageProducer(grpc_zap.DefaultMessageProducer),
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger.Logger)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			}),
			grpc_zap.UnaryServerInterceptor(logger.Logger, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.PayloadStreamServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			}),
			grpc_zap.StreamServerInterceptor(logger.Logger, opts...),
		),
	)
	pb.RegisterSatelitServer(grpcServer, s)

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// ShowVolume call GetVolume to Europa Backend
func (s *SatelitServer) ShowVolume(ctx context.Context, req *pb.ShowVolumeRequest) (*pb.ShowVolumeResponse, error) {
	volume, err := s.Europa.GetVolume(ctx, req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume: %+v", err)
	}

	return &pb.ShowVolumeResponse{
		Volume: volume.ToPb(),
	}, nil
}

// ListVolume call ListVolume to Europa Backend
func (s *SatelitServer) ListVolume(ctx context.Context, req *pb.ListVolumeRequest) (*pb.ListVolumeResponse, error) {
	volumes, err := s.Europa.ListVolume(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volumes: %+v", err)
	}

	var pvs []*pb.Volume
	for _, v := range volumes {
		pvs = append(pvs, v.ToPb())
	}

	return &pb.ListVolumeResponse{
		Volumes: pvs,
	}, nil
}

// AddVolume call CreateVolume to Europa backend
func (s *SatelitServer) AddVolume(ctx context.Context, req *pb.AddVolumeRequest) (*pb.AddVolumeResponse, error) {
	u, err := s.parseRequestUUID(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request name (need uuid): %+v", err)
	}

	volume, err := s.Europa.CreateVolume(ctx, u, int(req.CapacityGigabyte))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume: %+v", err)
	}

	return &pb.AddVolumeResponse{
		Volume: volume.ToPb(),
	}, nil
}

// AddVolumeImage call CreateVolumeImage to Europa backend
func (s *SatelitServer) AddVolumeImage(ctx context.Context, req *pb.AddVolumeImageRequest) (*pb.AddVolumeImageResponse, error) {
	u, err := s.parseRequestUUID(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request name (need uuid): %+v", err)
	}

	sourceImageID, err := s.parseRequestUUID(req.SourceImageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request source image id (need uuid): %+v", err)
	}

	v, err := s.Europa.CreateVolumeFromImage(ctx, u, int(req.CapacityGigabyte), sourceImageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume from image: %+v", err)
	}

	return &pb.AddVolumeImageResponse{
		Volume: v.ToPb(),
	}, nil
}

// AttachVolume call AttachVolume to Europa backend
func (s *SatelitServer) AttachVolume(ctx context.Context, req *pb.AttachVolumeRequest) (*pb.AttachVolumeResponse, error) {
	_, _, err := s.Europa.AttachVolumeTeleskop(ctx, req.Id, req.Hostname)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach volume: %+v", err)
	}

	return &pb.AttachVolumeResponse{}, nil
}

// DetachVolume call DetachVolume to Europa backend
func (s *SatelitServer) DetachVolume(ctx context.Context, req *pb.DetachVolumeRequest) (*pb.DetachVolumeResponse, error) {
	err := s.Europa.DetachVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to detach volume: %+v", err)
	}

	return &pb.DetachVolumeResponse{}, nil
}

// DeleteVolume call DeleteVolume to Europa backend
func (s *SatelitServer) DeleteVolume(ctx context.Context, req *pb.DeleteVolumeRequest) (*pb.DeleteVolumeResponse, error) {
	err := s.Europa.DeleteVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete volume: %+v", err)
	}

	return nil, nil
}

// parseRequestUUID return uuid.UUID from gRPC request string
func (s *SatelitServer) parseRequestUUID(reqUUID string) (uuid.UUID, error) {
	u := uuid.FromStringOrNil(reqUUID)
	if u == uuid.Nil {
		return uuid.Nil, fmt.Errorf("failed to parse uuid from string (uuid: %s)", reqUUID)
	}

	return u, nil
}

// ListImage retrieves all images
func (s *SatelitServer) ListImage(ctx context.Context, req *pb.ListImageRequest) (*pb.ListImageResponse, error) {
	images, err := s.Europa.ListImage()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieves images: %+v", err)
	}

	var pbImages []*pb.Image
	for _, image := range images {
		pbImages = append(pbImages, image.ToPb())
	}

	return &pb.ListImageResponse{
		Images: pbImages,
	}, nil
}

// UploadImage upload to europa backend
func (s *SatelitServer) UploadImage(stream pb.Satelit_UploadImageServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	m, err := s.receiveImage(stream, buf)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive image file: %+v", err)
	}
	logger.Logger.Debug(fmt.Sprintf("received image (name: %s)", m.name))

	// validate qcow2 image
	b := buf.Bytes()
	reader := bytes.NewReader(b)
	isQcow2, header := qcow2.Probe(reader)
	if isQcow2 == false {
		return status.Errorf(codes.InvalidArgument, "failed to validate qcow2 image")
	}

	// send to europa
	image, err := s.Europa.UploadImage(ctx, b, m.name, m.description, sanitizeImageSize(header.Size))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upload image to europa: %+v", err)
	}
	logger.Logger.Debug("uploaded image to europa")

	err = stream.SendAndClose(&pb.UploadImageResponse{Image: image.ToPb()})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to send and close: %+v", err)
	}
	logger.Logger.Debug("close UploadImage stream")

	// save to image info in database
	err = s.Datastore.PutImage(*image)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to put image: %+v", err)
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

	return m, nil
}

// DeleteImage delete image
func (s *SatelitServer) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.DeleteImageResponse, error) {
	imageID, err := s.parseRequestUUID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request image id (need uuid): %+v", err)
	}

	err = s.Europa.DeleteImage(ctx, imageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete image from europa: %+v", err)
	}

	return &pb.DeleteImageResponse{}, nil
}

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
	mac, err := net.ParseMAC(req.MacAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse mac address: %+v", err)
	}
	lease, err := s.IPAM.GetLease(ctx, mac)
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
	mac, err := net.ParseMAC(req.MacAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse mac address: %+v", err)
	}

	if err := s.IPAM.DeleteLease(ctx, mac); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete lease: %+v", err)
	}

	return &pb.DeleteLeaseResponse{}, nil
}

// AddVirtualMachine create virtual machine.
func (s *SatelitServer) AddVirtualMachine(ctx context.Context, req *pb.AddVirtualMachineRequest) (*pb.AddVirtualMachineResponse, error) {
	sourceImageID, err := s.parseRequestUUID(req.SourceImageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request source image id (need uuid): %+v", err)
	}

	u := uuid.NewV4()
	volume, err := s.Europa.CreateVolumeFromImage(ctx, u, int(req.RootVolumeGb), sourceImageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume from image: %+v", err)
	}

	_, deviceName, err := s.Europa.AttachVolumeTeleskop(ctx, volume.ID, req.HypervisorName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach volume: %+v", err)
	}

	vm, err := s.Ganymede.CreateVirtualMachine(ctx, req.Name, req.Vcpus, req.MemoryKib, deviceName, req.HypervisorName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create virtual machine: %+v", err)
	}

	return &pb.AddVirtualMachineResponse{
		Name: vm.Name,
		Uuid: vm.UUID.String(),
	}, nil
}

// StartVirtualMachine start virtual machine
func (s *SatelitServer) StartVirtualMachine(ctx context.Context, req *pb.StartVirtualMachineRequest) (*pb.StartVirtualMachineResponse, error) {
	vm, err := s.Datastore.GetVirtualMachine(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve virtual machine: %+v", err)
	}

	teleskopClient, err := teleskop.GetClient(vm.HypervisorName)
	if errors.Is(err, teleskop.ErrTeleskopAgentNotFound) {
		return nil, status.Errorf(codes.NotFound, "failed to retrieve teleskop client: %+v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve teleskop client: %+v", err)
	}

	resp, err := teleskopClient.StartVirtualMachine(ctx, &agentpb.StartVirtualMachineRequest{Uuid: req.Uuid})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start virtual machine: %+v", err)
	}

	return &pb.StartVirtualMachineResponse{
		Uuid: resp.Uuid,
		Name: resp.Name,
	}, nil
}
