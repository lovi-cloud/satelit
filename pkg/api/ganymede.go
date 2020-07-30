package api

import (
	"context"
	"errors"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/client/teleskop"
	agentpb "github.com/whywaita/teleskop/protoc/agent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddVirtualMachine create virtual machine.
func (s *SatelitServer) AddVirtualMachine(ctx context.Context, req *pb.AddVirtualMachineRequest) (*pb.AddVirtualMachineResponse, error) {
	sourceImageID, err := s.parseRequestUUID(req.SourceImageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req source image id (need uuid): %+v", err)
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

	vm, err := s.Ganymede.CreateVirtualMachine(ctx, req.Name, req.Vcpus, req.MemoryKib, deviceName, req.HypervisorName, volume.ID)
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
	vmID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req virtual machine id (need uuid): %+v", err)
	}

	vm, err := s.Datastore.GetVirtualMachine(vmID)
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

// ShowVirtualMachine retrieves virtual machine
func (s *SatelitServer) ShowVirtualMachine(ctx context.Context, req *pb.ShowVirtualMachineRequest) (*pb.ShowVirtualMachineResponse, error) {
	vmID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req virtual machine id (need uuid): %+v", err)
	}

	vm, err := s.Datastore.GetVirtualMachine(vmID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve virtual machine: %+v", err)
	}

	return &pb.ShowVirtualMachineResponse{
		VirtualMachine: vm.ToPb(),
	}, nil
}

// DeleteVirtualMachine delete virtual machine
func (s *SatelitServer) DeleteVirtualMachine(ctx context.Context, req *pb.DeleteVirtualMachineRequest) (*pb.DeleteVirtualMachineResponse, error) {
	vmID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req virtual machine id (need uuid): %+v", err)
	}

	vm, err := s.Datastore.GetVirtualMachine(vmID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve virtual machine: %+v", err)
	}

	err = s.Ganymede.DeleteVirtualMachine(ctx, vm.UUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete virtual machine: %+v", err)
	}

	err = s.Europa.DeleteVolume(ctx, vm.RootVolumeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete root volume: %+v", err)
	}

	return &pb.DeleteVirtualMachineResponse{}, nil
}

// CreateBridge is
func (s *SatelitServer) CreateBridge(ctx context.Context, req *pb.CreateBridgeRequest) (*pb.CreateBridgeResponse, error) {
	bridge, err := s.Ganymede.CreateBridge(ctx, uint32(req.VlanId))
	if err != nil {
		return nil, err
	}

	return &pb.CreateBridgeResponse{
		Bridge: bridge.ToPb(),
	}, nil
}

// GetBridge is
func (s *SatelitServer) GetBridge(ctx context.Context, req *pb.GetBridgeRequest) (*pb.GetBridgeResponse, error) {
	bridgeID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req bridge id (need uuid): %+v", err)
	}
	bridge, err := s.Ganymede.GetBridge(ctx, bridgeID)
	if err != nil {
		return nil, err
	}

	return &pb.GetBridgeResponse{
		Bridge: bridge.ToPb(),
	}, nil
}

// ListBridge is
func (s *SatelitServer) ListBridge(ctx context.Context, req *pb.ListBridgeRequest) (*pb.ListBridgeResponse, error) {
	bs, err := s.Ganymede.ListBridge(ctx)
	if err != nil {
		return nil, err
	}

	bridges := make([]*pb.Bridge, len(bs))
	for i, b := range bs {
		bridges[i] = b.ToPb()
	}

	return &pb.ListBridgeResponse{
		Bridges: bridges,
	}, nil
}

// DeleteBridge is
func (s *SatelitServer) DeleteBridge(ctx context.Context, req *pb.DeleteBridgeRequest) (*pb.DeleteBridgeResponse, error) {
	bridgeID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req bridge id (need uuid): %+v", err)
	}
	err = s.Ganymede.DeleteBridge(ctx, bridgeID)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteBridgeResponse{}, nil
}

// AttachInterface is
func (s *SatelitServer) AttachInterface(ctx context.Context, req *pb.AttachInterfaceRequest) (*pb.AttachInterfaceResponse, error) {
	vmID, err := s.parseRequestUUID(req.VirtualMachineId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req vm id: %+v", err)
	}
	fmt.Fprintf(os.Stderr, "bridge: %s\n", req.BridgeId)
	bridgeID, err := s.parseRequestUUID(req.BridgeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req bridge id: %+v", err)
	}
	leaseID, err := s.parseRequestUUID(req.LeaseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req bridge id: %+v", err)
	}

	attachment, err := s.Ganymede.AttachInterface(ctx, vmID, bridgeID, leaseID, int(req.Average), req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach interface: %+v", err)
	}

	return &pb.AttachInterfaceResponse{
		InterfaceAttachment: attachment.ToPb(),
	}, nil
}

// DetachInterface is
func (s *SatelitServer) DetachInterface(ctx context.Context, req *pb.DetachInterfaceRequest) (*pb.DetachInterfaceResponse, error) {
	attachmentID, err := s.parseRequestUUID(req.AtttachmentId)
	if err != nil {
		return nil, err
	}
	err = s.Ganymede.DetachInterface(ctx, attachmentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to detach interface: %+v", err)
	}

	return &pb.DetachInterfaceResponse{}, nil
}

// GetAttachment is
func (s *SatelitServer) GetAttachment(ctx context.Context, req *pb.GetAttachmentRequest) (*pb.GetAttachmentResponse, error) {
	attachmentID, err := s.parseRequestUUID(req.AttachmentId)
	if err != nil {
		return nil, err
	}
	attachent, err := s.Ganymede.GetAttachment(ctx, attachmentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get attachment: %+v", err)
	}

	return &pb.GetAttachmentResponse{
		InterfaceAttachment: attachent.ToPb(),
	}, nil
}

// ListAttachment is
func (s *SatelitServer) ListAttachment(ctx context.Context, req *pb.ListAttachmentRequest) (*pb.ListAttachmentResponse, error) {
	as, err := s.Ganymede.ListAttachment(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get attachment list: %+v", err)
	}

	attachments := make([]*pb.InterfaceAttachment, len(as))
	for i, a := range attachments {
		attachments[i] = a
	}

	return &pb.ListAttachmentResponse{
		InterfaceAttachments: attachments,
	}, nil
}
