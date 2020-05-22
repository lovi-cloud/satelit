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

	u := uuid.FromStringOrNil(testUUID)

	_, err = m.CreateVolumeRaw(ctx, u, testCapacity)
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
	volume := vs[0]

	if volume.ID != testUUID {
		t.Error("Unexpected volume id")
	}

	_, _, err = m.AttachVolume(ctx, volume.ID, testHostname)
	if err != nil {
		t.Error(err)
	}

	v, err := m.GetVolume(ctx, volume.ID)
	if err != nil {
		t.Error("Can't get volume")
	}

	if v.Attached != true || v.HostName != testHostname {
		t.Error("Unexpected volume attachment info after AttachVolume")
	}

	err = m.DetachVolume(ctx, volume.ID)
	if err != nil {
		t.Error(err)
	}

	v, err = m.GetVolume(ctx, volume.ID)
	if err != nil {
		t.Error("Can't get volume")
	}

	if v.Attached != false || v.HostName != "" {
		t.Error("Unexpected volume attachment info after DetachVolume")
	}

}
