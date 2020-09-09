package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/internal/mysql/types"

	"github.com/whywaita/satelit/internal/testutils"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
)

const (
	testVirtualMachineID = "7b3135c4-a1ed-4749-aa13-53de2fef678e"
	testRootVolumeID     = "22fe4009-aba7-48ce-9b54-79520ffa2b4a"
)

func TestMySQL_GetVirtualMachine(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v\n", err)
	}

	err = testDatastore.PutVolume(europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "hv000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	err = testDatastore.PutVirtualMachine(ganymede.VirtualMachine{
		UUID:           uuid.FromStringOrNil(testVirtualMachineID),
		Name:           "test000",
		Vcpus:          1,
		MemoryKiB:      2 * 1024 * 1024,
		HypervisorName: "hv000",
		RootVolumeID:   testRootVolumeID,
		ReadBytesSec:   100 * 1000 * 1000,
		WriteBytesSec:  200 * 1000 * 1000,
		ReadIOPSSec:    10000,
		WriteIOPSSec:   5000,
	})
	if err != nil {
		t.Fatalf("failed to put virtual machine: %+v\n", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.VirtualMachine
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testVirtualMachineID),
			want: &ganymede.VirtualMachine{
				UUID:           uuid.FromStringOrNil(testVirtualMachineID),
				Name:           "test000",
				Vcpus:          1,
				MemoryKiB:      2 * 1024 * 1024,
				HypervisorName: "hv000",
				RootVolumeID:   testRootVolumeID,
				RootVolumeGB:   20,
				ReadBytesSec:   100 * 1000 * 1000,
				WriteBytesSec:  200 * 1000 * 1000,
				ReadIOPSSec:    10000,
				WriteIOPSSec:   5000,
				SourceImageID:  testImage.UUID,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetVirtualMachine(test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_PutVirtualMachine(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v\n", err)
	}

	err = testDatastore.PutVolume(europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "hv000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	testCorePairUUIDs := []uuid.UUID{
		uuid.FromStringOrNil("9cf11645-ec85-4607-b638-cd592819bbae"),
		uuid.FromStringOrNil("25b403a9-cdd7-4176-8d44-c922220bdcb8"),
		uuid.FromStringOrNil("2cc61359-8912-4187-aadc-8692574b1b52"),
		uuid.FromStringOrNil("e77523a3-fef0-4864-b24f-4f9579a65eed"),
	}

	testVirtualMachineID2 := "575abf68-c10e-4801-84df-8bc1bed82ff2"

	_, err = setFakePinnedCore(testDatastore, testCorePairUUIDs)
	if err != nil {
		t.Fatalf("failed to set cpu pinned core: %+v", err)
	}

	cpg, err := testDatastore.GetCPUPinningGroup(context.Background(), uuid.FromStringOrNil(testCPUPinningGroupUUID))
	if err != nil {
		t.Fatalf("failed to retrieve cpu pinning group: %+v", err)
	}

	tests := []struct {
		input ganymede.VirtualMachine
		want  *ganymede.VirtualMachine
		err   bool
	}{
		{
			input: ganymede.VirtualMachine{
				UUID:              uuid.FromStringOrNil(testVirtualMachineID),
				Name:              "test000",
				Vcpus:             1,
				MemoryKiB:         2 * 1024 * 1024,
				HypervisorName:    "hv000",
				RootVolumeID:      testRootVolumeID,
				RootVolumeGB:      20,
				ReadBytesSec:      100 * 1000 * 1000,
				WriteBytesSec:     200 * 1000 * 1000,
				ReadIOPSSec:       10000,
				WriteIOPSSec:      5000,
				SourceImageID:     testImage.UUID,
				CPUPinningGroupID: cpg.UUID,
			},
			want: &ganymede.VirtualMachine{
				UUID:              uuid.FromStringOrNil(testVirtualMachineID),
				Name:              "test000",
				Vcpus:             1,
				MemoryKiB:         2 * 1024 * 1024,
				HypervisorName:    "hv000",
				RootVolumeID:      testRootVolumeID,
				RootVolumeGB:      20,
				ReadBytesSec:      100 * 1000 * 1000,
				WriteBytesSec:     200 * 1000 * 1000,
				ReadIOPSSec:       10000,
				WriteIOPSSec:      5000,
				SourceImageID:     testImage.UUID,
				CPUPinningGroupID: cpg.UUID,
			},
			err: false,
		}, {
			input: ganymede.VirtualMachine{
				UUID:           uuid.FromStringOrNil(testVirtualMachineID2),
				Name:           "test000-no-cpu-pinning-group",
				Vcpus:          1,
				MemoryKiB:      2 * 1024 * 1024,
				HypervisorName: "hv000",
				RootVolumeID:   testRootVolumeID,
				RootVolumeGB:   20,
				ReadBytesSec:   100 * 1000 * 1000,
				WriteBytesSec:  200 * 1000 * 1000,
				ReadIOPSSec:    10000,
				WriteIOPSSec:   5000,
				SourceImageID:  testImage.UUID,
			},
			want: &ganymede.VirtualMachine{
				UUID:           uuid.FromStringOrNil(testVirtualMachineID2),
				Name:           "test000-no-cpu-pinning-group",
				Vcpus:          1,
				MemoryKiB:      2 * 1024 * 1024,
				HypervisorName: "hv000",
				RootVolumeID:   testRootVolumeID,
				RootVolumeGB:   20,
				ReadBytesSec:   100 * 1000 * 1000,
				WriteBytesSec:  200 * 1000 * 1000,
				ReadIOPSSec:    10000,
				WriteIOPSSec:   5000,
				SourceImageID:  testImage.UUID,
			},
			err: false,
		},
	}
	for _, test := range tests {
		err := testDatastore.PutVirtualMachine(test.input)
		if err != nil {
			t.Fatalf("failed to create subnet: %+v", err)
		}
		got, err := getVirtualMachineFromSQL(testDB, test.input.UUID)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_ListVirtualMachine(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	//testDB, _ := testutils.GetTestDB()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v\n", err)
	}

	err = testDatastore.PutVolume(europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "hv000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	err = testDatastore.PutVirtualMachine(ganymede.VirtualMachine{
		UUID:           uuid.FromStringOrNil(testVirtualMachineID),
		Name:           "test000",
		Vcpus:          1,
		MemoryKiB:      2 * 1024 * 1024,
		HypervisorName: "hv000",
		RootVolumeID:   testRootVolumeID,
		RootVolumeGB:   20,
		ReadBytesSec:   100 * 1000 * 1000,
		WriteBytesSec:  200 * 1000 * 1000,
		ReadIOPSSec:    10000,
		WriteIOPSSec:   5000,
		SourceImageID:  testImage.UUID,
	})
	if err != nil {
		t.Fatalf("failed to put virtual machine: %+v\n", err)
	}

	tests := []struct {
		input interface{}
		want  []ganymede.VirtualMachine
		err   bool
	}{
		{
			input: nil,
			want: []ganymede.VirtualMachine{
				{
					UUID:           uuid.FromStringOrNil(testVirtualMachineID),
					Name:           "test000",
					Vcpus:          1,
					MemoryKiB:      2 * 1024 * 1024,
					HypervisorName: "hv000",
					RootVolumeID:   testRootVolumeID,
					RootVolumeGB:   20,
					ReadBytesSec:   100 * 1000 * 1000,
					WriteBytesSec:  200 * 1000 * 1000,
					ReadIOPSSec:    10000,
					WriteIOPSSec:   5000,
					SourceImageID:  testImage.UUID,
				},
			},
		},
	}

	for _, test := range tests {
		got, err := testDatastore.ListVirtualMachine()
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_DeleteVirtualMachine(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v\n", err)
	}

	err = testDatastore.PutVolume(europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "hv000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	err = testDatastore.PutVirtualMachine(ganymede.VirtualMachine{
		UUID:           uuid.FromStringOrNil(testVirtualMachineID),
		Name:           "test000",
		Vcpus:          1,
		MemoryKiB:      2 * 1024 * 1024,
		HypervisorName: "hv000",
		RootVolumeID:   testRootVolumeID,
		RootVolumeGB:   20,
		ReadBytesSec:   100 * 1000 * 1000,
		WriteBytesSec:  200 * 1000 * 1000,
		ReadIOPSSec:    10000,
		WriteIOPSSec:   5000,
		SourceImageID:  testImage.UUID,
	})
	if err != nil {
		t.Fatalf("failed to put virtual machine: %+v\n", err)
	}

	tests := []struct {
		input uuid.UUID
		want  *ganymede.VirtualMachine
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testVirtualMachineID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DeleteVirtualMachine(test.input)
		if err != nil {
			t.Fatalf("failed to delete virtual machine: %+v", err)
		}
		got, err := getVirtualMachineFromSQL(testDB, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_GetHostnameByAddress(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	bridge, err := testDatastore.CreateBridge(context.Background(), ganymede.Bridge{
		UUID:   uuid.FromStringOrNil(testBridgeID),
		VLANID: 1000,
		Name:   "testbr1000",
	})
	if err != nil {
		t.Fatalf("failed to create bridge: %+v", err)
	}

	err = testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v\n", err)
	}

	err = testDatastore.PutVolume(europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "hv000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	err = testDatastore.PutVirtualMachine(ganymede.VirtualMachine{
		UUID:           uuid.FromStringOrNil(testVirtualMachineID),
		Name:           "test000",
		Vcpus:          1,
		MemoryKiB:      2 * 1024 * 1024,
		HypervisorName: "hv000",
		RootVolumeID:   testRootVolumeID,
		ReadBytesSec:   100 * 1000 * 1000,
		WriteBytesSec:  200 * 1000 * 1000,
		ReadIOPSSec:    10000,
		WriteIOPSSec:   5000,
	})
	if err != nil {
		t.Fatalf("failed to put virtual machine: %+v\n", err)
	}

	_, err = testDatastore.CreateSubnet(context.Background(), testSubnet)
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	_, err = testDatastore.CreateAddress(context.Background(), testAddress)
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	lease, err := testDatastore.CreateLease(context.Background(), testLease)
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	_, err = testDatastore.AttachInterface(context.Background(), ganymede.InterfaceAttachment{
		UUID:             uuid.FromStringOrNil(testAttachmentID),
		VirtualMachineID: uuid.FromStringOrNil(testVirtualMachineID),
		BridgeID:         bridge.UUID,
		Average:          1 * 1024 * 1024,
		Name:             "vnet000",
		LeaseID:          lease.UUID,
	})
	if err != nil {
		t.Fatalf("failed to create test attachment: %+v", err)
	}

	tests := []struct {
		input types.IP
		want  string
		err   bool
	}{
		{
			input: testAddress.IP,
			want:  "test000",
			err:   false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetHostnameByAddress(test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func getVirtualMachineFromSQL(testDB *sqlx.DB, vmID uuid.UUID) (*ganymede.VirtualMachine, error) {
	query := `SELECT uuid, name, vcpus, memory_kib, hypervisor_name, root_volume_id, volume.capacity_gb, read_bytes_sec, write_bytes_sec, read_iops_sec, write_iops_sec, volume.base_image_id, cpu_pinning_group_id FROM virtual_machine JOIN volume ON virtual_machine.root_volume_id = volume.id WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var v ganymede.VirtualMachine
	err = stmt.Get(&v, vmID)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
