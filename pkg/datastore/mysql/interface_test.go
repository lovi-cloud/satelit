package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/jmoiron/sqlx"

	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/satelit/internal/testutils"
	"github.com/lovi-cloud/satelit/pkg/europa"
	"github.com/lovi-cloud/satelit/pkg/ganymede"
)

const (
	testAttachmentID = "c0aff0e7-db1b-4e10-96ff-8a701de290c5"
)

func TestMySQL_AttachInterface(t *testing.T) {
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

	err = testDatastore.PutVolume(context.Background(), europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
		BackendName: testVolume.BackendName,
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

	tests := []struct {
		input ganymede.InterfaceAttachment
		want  *ganymede.InterfaceAttachment
		err   bool
	}{
		{
			input: ganymede.InterfaceAttachment{
				UUID:             uuid.FromStringOrNil(testAttachmentID),
				VirtualMachineID: uuid.FromStringOrNil(testVirtualMachineID),
				BridgeID:         bridge.UUID,
				Average:          1 * 1024 * 1024,
				Name:             "vnet000",
				LeaseID:          lease.UUID,
			},
			want: &ganymede.InterfaceAttachment{
				UUID:             uuid.FromStringOrNil(testAttachmentID),
				VirtualMachineID: uuid.FromStringOrNil(testVirtualMachineID),
				BridgeID:         bridge.UUID,
				Average:          1 * 1024 * 1024,
				Name:             "vnet000",
				LeaseID:          lease.UUID,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.AttachInterface(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_DetachInterface(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

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

	err = testDatastore.PutVolume(context.Background(), europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
		BackendName: testVolume.BackendName,
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
		input uuid.UUID
		want  *ganymede.InterfaceAttachment
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testAttachmentID),
			want:  nil,
			err:   true,
		},
	}
	for _, test := range tests {
		err := testDatastore.DetachInterface(context.Background(), test.input)
		if err != nil {
			t.Fatalf("failed to detach interface: %+v", err)
		}
		got, err := getAttachmentFromSQL(testDB, test.input)
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

func TestMySQL_GetAttachment(t *testing.T) {
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

	err = testDatastore.PutVolume(context.Background(), europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
		BackendName: testVolume.BackendName,
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
		input uuid.UUID
		want  *ganymede.InterfaceAttachment
		err   bool
	}{
		{
			input: uuid.FromStringOrNil(testAttachmentID),
			want: &ganymede.InterfaceAttachment{
				UUID:             uuid.FromStringOrNil(testAttachmentID),
				VirtualMachineID: uuid.FromStringOrNil(testVirtualMachineID),
				BridgeID:         bridge.UUID,
				Average:          1 * 1024 * 1024,
				Name:             "vnet000",
				LeaseID:          lease.UUID,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.GetAttachment(context.Background(), test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if got != nil {
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestMySQL_ListAttachment(t *testing.T) {
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

	err = testDatastore.PutVolume(context.Background(), europa.Volume{
		ID:          testRootVolumeID,
		Attached:    false,
		HostName:    "",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
		BackendName: testVolume.BackendName,
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
		want []ganymede.InterfaceAttachment
		err  bool
	}{
		{
			want: []ganymede.InterfaceAttachment{
				{
					UUID:             uuid.FromStringOrNil(testAttachmentID),
					VirtualMachineID: uuid.FromStringOrNil(testVirtualMachineID),
					BridgeID:         bridge.UUID,
					Average:          1 * 1024 * 1024,
					Name:             "vnet000",
					LeaseID:          lease.UUID,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := testDatastore.ListAttachment(context.Background())
		if !test.err && err != nil {
			t.Fatalf("should not be error but: %+v", err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error but not:")
		}
		if got != nil {
			for i := 0; i < len(got); i++ {
				got[i].CreatedAt = time.Time{}
				got[i].UpdatedAt = time.Time{}
			}
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func getAttachmentFromSQL(testDB *sqlx.DB, attachmentID uuid.UUID) (*ganymede.InterfaceAttachment, error) {
	query := `SELECT uuid, virtual_machine_id, bridge_id, average, name, lease_id, created_at, updated_at FROM interface_attachment WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var a ganymede.InterfaceAttachment
	err = stmt.Get(&a, attachmentID)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
