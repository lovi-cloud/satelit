package api

import (
	"context"
	"errors"

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
	vmID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request virtual machine id (need uuid): %+v", err)
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
