package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/whywaita/satelit/pkg/ganymede"

	"github.com/whywaita/satelit/internal/logger"

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
	defer func() {
		volumeID := volume.ID
		if err != nil {
			if err := s.Europa.DeleteVolume(ctx, volumeID); err != nil {
				logger.Logger.Warn(fmt.Sprintf("failed to DeleteVolume: %v", err))
			}
		}
	}()

	_, deviceName, err := s.Europa.AttachVolumeTeleskop(ctx, volume.ID, req.HypervisorName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach volume: %+v", err)
	}
	defer func() {
		volumeID := volume.ID
		if err != nil {
			if err := s.Europa.DetachVolume(ctx, volumeID); err != nil {
				logger.Logger.Warn(fmt.Sprintf("failed to DetachVolume: %v", err))
			}
		}
	}()

	vm, err := s.Ganymede.CreateVirtualMachine(ctx, req.Name, req.Vcpus, req.MemoryKib, deviceName, req.HypervisorName, volume.ID, req.ReadBytesSec, req.WriteBytesSec, req.ReadIopsSec, req.WriteIopsSec, req.PinningGroupName)
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

	cpgName := ""
	if vm.CPUPinningGroupID != uuid.Nil {
		cpg, err := s.Datastore.GetCPUPinningGroup(ctx, vm.CPUPinningGroupID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to retrieve cpu pinning group: %+v", err)
		}

		cpgName = cpg.Name
	}

	return &pb.ShowVirtualMachineResponse{
		VirtualMachine: vm.ToPb(cpgName),
	}, nil
}

// ListVirtualMachine retrieve all virtual machine
func (s *SatelitServer) ListVirtualMachine(ctx context.Context, req *pb.ListVirtualMachineRequest) (*pb.ListVirtualMachineResponse, error) {
	vms, err := s.Datastore.ListVirtualMachine()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve all virtual machine from datastore: %+v", err)
	}

	var pbvms []*pb.VirtualMachine
	for _, vm := range vms {
		cpgName := ""
		if vm.CPUPinningGroupID != uuid.Nil {
			cpg, err := s.Datastore.GetCPUPinningGroup(ctx, vm.CPUPinningGroupID)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to retrieve cpu pinning group from datastore: %+v", err)
			}

			cpgName = cpg.Name
		}
		pbvms = append(pbvms, vm.ToPb(cpgName))
	}

	return &pb.ListVirtualMachineResponse{
		VirtualMachines: pbvms,
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

	client, err := teleskop.GetClient(vm.HypervisorName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get teleskop client: %+v", err)
	}
	resp, err := client.GetVirtualMachineState(ctx, &agentpb.GetVirtualMachineStateRequest{Uuid: vmID.String()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve virtual machine state: %+v", err)
	}
	if resp.State.State.Type() != agentpb.VirtualMachineState_SHUTOFF.Type() {
		return nil, status.Errorf(codes.InvalidArgument, "%s is not shutdown. please `virsh destroy` before delete", vmID.String())
	}

	if err := s.Ganymede.DeleteVirtualMachine(ctx, vm.UUID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete virtual machine: %+v", err)
	}

	if err := s.Europa.DetachVolume(ctx, vm.RootVolumeID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to detach root volume: %+v", err)
	}
	if err := s.Europa.DeleteVolume(ctx, vm.RootVolumeID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete root volume: %+v", err)
	}

	return &pb.DeleteVirtualMachineResponse{}, nil
}

// CreateBridge is
func (s *SatelitServer) CreateBridge(ctx context.Context, req *pb.CreateBridgeRequest) (*pb.CreateBridgeResponse, error) {
	bridge, err := s.Ganymede.CreateBridge(ctx, req.Name, uint32(req.VlanId))
	if err != nil {
		return nil, err
	}

	return &pb.CreateBridgeResponse{
		Bridge: bridge.ToPb(),
	}, nil
}

// CreateInternalBridge is
func (s *SatelitServer) CreateInternalBridge(ctx context.Context, req *pb.CreateInternalBridgeRequest) (*pb.CreateInternalBridgeResponse, error) {
	bridge, err := s.Ganymede.CreateInternalBridge(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.CreateInternalBridgeResponse{
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
	bridgeID, err := s.parseRequestUUID(req.BridgeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req bridge id: %+v", err)
	}
	leaseID, err := s.parseRequestUUID(req.LeaseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req lease id: %+v", err)
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
	for i, a := range as {
		attachments[i] = a.ToPb()
	}

	return &pb.ListAttachmentResponse{
		InterfaceAttachments: attachments,
	}, nil
}

// AddCPUPinningGroup add cpu pinning group
// use same group's cpu cores if virtual machine joined a same cpu pinning group
func (s *SatelitServer) AddCPUPinningGroup(ctx context.Context, req *pb.AddCPUPinningGroupRequest) (*pb.AddCPUPinningGroupResponse, error) {
	div := req.CountOfCore % 2
	if div != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "count_of_core must be a multiple of two (physical and logical core)")
	}
	hv, err := s.Datastore.GetHypervisorByHostname(ctx, req.HypervisorName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve hypervisor: %+v", err)
	}

	u := uuid.NewV4()
	cpg := ganymede.CPUPinningGroup{
		UUID:         u,
		Name:         req.Name,
		CountCore:    int(req.CountOfCore),
		HypervisorID: hv.ID,
	}
	if err := s.Datastore.PutCPUPinningGroup(ctx, cpg); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to put CPU Pinning Group: %+v", err)
	}

	numRequestCorePair := req.CountOfCore / 2
	_, err = s.Scheduler.PopCorePair(ctx, hv.ID, int(numRequestCorePair), cpg.UUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to allocate pinning cpu pair: %+v", err)
	}
	defer func() {
		if err != nil {
			if err := s.Scheduler.PushCorePair(ctx, cpg.UUID); err != nil {
				logger.Logger.Warn(fmt.Sprintf("failed to push core pair: %+v", err))
			}
		}
	}()

	return &pb.AddCPUPinningGroupResponse{CpuPinningGroup: &pb.CPUPinningGroup{
		Uuid:        u.String(),
		Name:        req.Name,
		CountOfCore: req.CountOfCore,
	}}, nil
}

// ShowCPUPinningGroup retrieve cpu pinning group
func (s *SatelitServer) ShowCPUPinningGroup(ctx context.Context, req *pb.ShowCPUPinningGroupRequest) (*pb.ShowCPUPinningGroupResponse, error) {
	cpuPinningGroupID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req cpu pinning group id (need uuid): %+v", err)
	}

	cpg, err := s.Datastore.GetCPUPinningGroup(ctx, cpuPinningGroupID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cpu pinning group: %+v", err)
	}

	return &pb.ShowCPUPinningGroupResponse{
		CpuPinningGroup: cpg.ToPb(),
	}, nil
}

// DeleteCPUPinningGroup delete cpu pinning group
func (s *SatelitServer) DeleteCPUPinningGroup(ctx context.Context, req *pb.DeleteCPUPinningGroupRequest) (*pb.DeleteCPUPinningGroupResponse, error) {
	cpuPinningGroupID, err := s.parseRequestUUID(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse req cpu pinning group id (need uuid): %+v", err)
	}

	// TODO: check related pinned cpu from hypervisor (return InvalidArgument if exist)
	cpg, err := s.Datastore.GetCPUPinningGroup(ctx, cpuPinningGroupID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cpu pinning group: %+v", err)
	}

	if err := s.Scheduler.PushCorePair(ctx, cpg.UUID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to free cpu core pair: %+v", err)
	}

	if err := s.Datastore.DeleteCPUPinningGroup(ctx, cpuPinningGroupID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete cpu pinning group: %+v", err)
	}

	return &pb.DeleteCPUPinningGroupResponse{}, nil
}
