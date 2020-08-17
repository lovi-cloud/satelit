package testutils

import (
	"context"
	"fmt"
	"net"
	"os"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"

	"github.com/whywaita/teleskop/protoc/agent"
)

type dummyTeleskop struct{}

// NewDummyTeleskop is
func NewDummyTeleskop() (string, func(), error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer()
	agent.RegisterAgentServer(server, &dummyTeleskop{})

	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return lis.Addr().String(), server.Stop, nil
}

func (d dummyTeleskop) GetISCSIQualifiedName(ctx context.Context, req *agent.GetISCSIQualifiedNameRequest) (*agent.GetISCSIQualifiedNameResponse, error) {
	return &agent.GetISCSIQualifiedNameResponse{
		Iqn: "sample-iqn",
	}, nil
}

func (d dummyTeleskop) GetIPTables(ctx context.Context, req *agent.GetIPTablesRequest) (*agent.GetIPTablesResponse, error) {
	return &agent.GetIPTablesResponse{}, nil
}

func (d dummyTeleskop) SetupDefaultSecurityGroup(ctx context.Context, req *agent.SetupDefaultSecurityGroupRequest) (*agent.SetupDefaultSecurityGroupResponse, error) {
	return &agent.SetupDefaultSecurityGroupResponse{}, nil
}

func (d dummyTeleskop) AddSecurityGroup(ctx context.Context, req *agent.AddSecurityGroupRequest) (*agent.AddSecurityGroupResponse, error) {
	return &agent.AddSecurityGroupResponse{}, nil
}

func (d dummyTeleskop) AddBridge(ctx context.Context, req *agent.AddBridgeRequest) (*agent.AddBridgeResponse, error) {
	return &agent.AddBridgeResponse{}, nil
}

func (d dummyTeleskop) AddVLANInterface(ctx context.Context, req *agent.AddVLANInterfaceRequest) (*agent.AddVLANInterfaceResponse, error) {
	return &agent.AddVLANInterfaceResponse{}, nil
}

func (d dummyTeleskop) AddInterfaceToBridge(ctx context.Context, req *agent.AddInterfaceToBridgeRequest) (*agent.AddInterfaceToBridgeResponse, error) {
	return &agent.AddInterfaceToBridgeResponse{}, nil
}

func (d dummyTeleskop) AddVirtualMachine(ctx context.Context, req *agent.AddVirtualMachineRequest) (*agent.AddVirtualMachineResponse, error) {
	return &agent.AddVirtualMachineResponse{
		Uuid: uuid.NewV4().String(),
		Name: req.Name,
	}, nil
}

func (d dummyTeleskop) ConnectBlockDevice(ctx context.Context, req *agent.ConnectBlockDeviceRequest) (*agent.ConnectBlockDeviceResponse, error) {
	return &agent.ConnectBlockDeviceResponse{
		DeviceName: "/dev/dm-0",
	}, nil
}

func (d dummyTeleskop) StartVirtualMachine(ctx context.Context, req *agent.StartVirtualMachineRequest) (*agent.StartVirtualMachineResponse, error) {
	return &agent.StartVirtualMachineResponse{
		Uuid: req.Uuid,
		Name: req.Uuid,
	}, nil
}

func (d dummyTeleskop) AttachBlockDevice(ctx context.Context, req *agent.AttachBlockDeviceRequest) (*agent.AttachBlockDeviceResponse, error) {
	return &agent.AttachBlockDeviceResponse{
		Uuid: req.Uuid,
		Name: req.Uuid,
	}, nil
}

func (d dummyTeleskop) AttachInterface(ctx context.Context, req *agent.AttachInterfaceRequest) (*agent.AttachInterfaceResponse, error) {
	return &agent.AttachInterfaceResponse{
		Uuid: req.Uuid,
		Name: req.Uuid,
	}, nil
}

func (d dummyTeleskop) DeleteBridge(ctx context.Context, req *agent.DeleteBridgeRequest) (*agent.DeleteBridgeResponse, error) {
	return &agent.DeleteBridgeResponse{}, nil
}

func (d dummyTeleskop) DeleteVLANInterface(ctx context.Context, req *agent.DeleteVLANInterfaceRequest) (*agent.DeleteVLANInterfaceResponse, error) {
	return &agent.DeleteVLANInterfaceResponse{}, nil
}

func (d dummyTeleskop) DeleteInterfaceFromBridge(ctx context.Context, req *agent.DeleteInterfaceFromBridgeRequest) (*agent.DeleteInterfaceFromBridgeResponse, error) {
	return &agent.DeleteInterfaceFromBridgeResponse{}, nil
}

func (d dummyTeleskop) DeleteVirtualMachine(ctx context.Context, req *agent.DeleteVirtualMachineRequest) (*agent.DeleteVirtualMachineResponse, error) {
	return &agent.DeleteVirtualMachineResponse{}, nil
}

func (d dummyTeleskop) DisconnectBlockDevice(ctx context.Context, req *agent.DisconnectBlockDeviceRequest) (*agent.DisconnectBlockDeviceResponse, error) {
	return &agent.DisconnectBlockDeviceResponse{}, nil
}

func (d dummyTeleskop) StopVirtualMachine(ctx context.Context, req *agent.StopVirtualMachineRequest) (*agent.StopVirtualMachineResponse, error) {
	return &agent.StopVirtualMachineResponse{}, nil
}

func (d dummyTeleskop) DetachBlockDevice(ctx context.Context, req *agent.DetachBlockDeviceRequest) (*agent.DetachBlockDeviceResponse, error) {
	return &agent.DetachBlockDeviceResponse{}, nil
}

func (d dummyTeleskop) DetachInterface(ctx context.Context, req *agent.DetachInterfaceRequest) (*agent.DetachInterfaceResponse, error) {
	return &agent.DetachInterfaceResponse{}, nil
}

func (d dummyTeleskop) GetVirtualMachineState(ctx context.Context, req *agent.GetVirtualMachineStateRequest) (*agent.GetVirtualMachineStateResponse, error) {
	return &agent.GetVirtualMachineStateResponse{
		State: &agent.VirtualMachineState{
			Uuid:  req.Uuid,
			Name:  req.Uuid,
			State: agent.VirtualMachineState_RUNNING,
		},
	}, nil
}

func (d dummyTeleskop) ListVirtualMachineState(ctx context.Context, req *agent.ListVirtualMachineStateRequest) (*agent.ListVirtualMachineStateResponse, error) {
	return &agent.ListVirtualMachineStateResponse{
		States: []*agent.VirtualMachineState{},
	}, nil
}
