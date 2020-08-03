package mysql_test

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

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
		HostName:    "dorad000",
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
			t.Fatalf("want %q, but %q, dirr %q:", test.want, got, diff)
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
		HostName:    "dorad000",
		CapacityGB:  20,
		BaseImageID: testImage.UUID,
		HostLUNID:   0,
	})
	if err != nil {
		t.Fatalf("failed to put volume: %+v\n", err)
	}

	tests := []struct {
		input ganymede.VirtualMachine
		want  *ganymede.VirtualMachine
		err   bool
	}{
		{
			input: ganymede.VirtualMachine{
				UUID:           uuid.FromStringOrNil(testVirtualMachineID),
				Name:           "test000",
				Vcpus:          1,
				MemoryKiB:      2 * 1024 * 1024,
				HypervisorName: "hv000",
				RootVolumeID:   testRootVolumeID,
			},
			want: &ganymede.VirtualMachine{
				UUID:           uuid.FromStringOrNil(testVirtualMachineID),
				Name:           "test000",
				Vcpus:          1,
				MemoryKiB:      2 * 1024 * 1024,
				HypervisorName: "hv000",
				RootVolumeID:   testRootVolumeID,
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
			t.Fatalf("want %q, but %q, dirr %q:", test.want, got, diff)
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
		HostName:    "dorad000",
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
			t.Fatalf("want %q, but %q, dirr %q:", test.want, got, diff)
		}
	}
}

func getVirtualMachineFromSQL(testDB *sqlx.DB, vmID uuid.UUID) (*ganymede.VirtualMachine, error) {
	query := `SELECT uuid, name, vcpus, memory_kib, hypervisor_name, root_volume_id FROM virtual_machine WHERE uuid = ?`
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