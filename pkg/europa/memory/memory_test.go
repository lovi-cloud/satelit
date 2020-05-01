package memory

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
)

const (
	testUUID     = "a52e6253-4324-4558-ab96-04b566b8de69"
	testCapacity = 1
	testHostname = "XXX"
)

func TestMemoryVolumeOperation(t *testing.T) {
	ctx := context.Background()
	m, err := New()
	if err != nil {
		t.Error(err)
	}

	_, err = m.CreateVolume(ctx, uuid.FromStringOrNil(testUUID), testCapacity)
	if err != nil {
		t.Error(err)
	}

	_, ok := m.Volumes[testUUID]
	if ok != true {
		t.Error("not created volume")
	}

	vs, err := m.ListVolume(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(vs) != 1 {
		t.Error("Unexpected num of volumes")
	}
	if vs[0].ID != testUUID {
		t.Error("Unexpected volume id")
	}

	_, err = m.AttachVolume(ctx, uuid.FromStringOrNil(testUUID), testHostname)
	if err != nil {
		t.Error(err)
	}
}
