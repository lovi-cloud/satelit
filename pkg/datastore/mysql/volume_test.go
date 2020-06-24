package mysql_test

import (
	"database/sql"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/whywaita/satelit/internal/testutils"
	"github.com/whywaita/satelit/pkg/europa"
)

var testVolume = europa.Volume{
	ID:          "TEST_VOLUME_ID_IS_NOT_UUID",
	CapacityGB:  8,
	BaseImageID: uuid.FromStringOrNil(testUUID),
}

func TestMySQL_PutVolume(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	// first write, use INSERT
	if err := testDatastore.PutVolume(testVolume); err != nil {
		t.Fatalf("failed to put volume (INSERT): %s", err)
	}

	v1, err := getVolumeFromSQL(testDB, testVolume.ID)
	if err != nil {
		t.Fatalf("failed to get volume from sql: %s", err)
	}

	ok, values := testutils.CompareStruct(testVolume, v1)
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}

	testVolume.Attached = true
	testVolume.HostName = "testHost"
	testVolume.HostLUNID = 1
	// second write, use UPDATE
	if err = testDatastore.PutVolume(testVolume); err != nil {
		t.Fatalf("failed to put volume (UPDATE):%s", err)
	}

	v2, err := getVolumeFromSQL(testDB, testVolume.ID)
	if err != nil {
		t.Fatalf("failed to get volume from sql: %s", err)
	}

	want := &europa.Volume{
		ID:          testVolume.ID,
		Attached:    true,
		HostName:    "testHost",
		CapacityGB:  8,
		BaseImageID: uuid.FromStringOrNil(testUUID),
		HostLUNID:   1,
	}
	ok, values = testutils.CompareStruct(want, v2)
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}
}

func TestMySQL_GetVolume(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	if err := testDatastore.PutVolume(testVolume); err != nil {
		t.Fatalf("failed to put volume (INSERT): %s", err)
	}

	volume, err := testDatastore.GetVolume(testVolume.ID)
	if err != nil {
		t.Fatalf("failed to get volume: %s", err)
	}

	want := &testVolume
	ok, values := testutils.CompareStruct(want, volume)
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}
}

func TestMySQL_DeleteVolume(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	if err := testDatastore.PutVolume(testVolume); err != nil {
		t.Fatalf("failed to put volume (INSERT): %s", err)
	}
	if err := testDatastore.DeleteVolume(testVolume.ID); err != nil {
		t.Fatalf("failed to delete volume: %s", err)
	}

	_, err := getVolumeFromSQL(testDB, testVolume.ID)
	if err == nil {
		t.Fatalf("return volume after DeleteVolume")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("failed to get volume from SQL: %s", err)
	}
}

func getVolumeFromSQL(testDB *sqlx.DB, volumeID string) (*europa.Volume, error) {
	query := fmt.Sprintf(`SELECT id, attached, hostname, capacity_gb, base_image_id, host_lun_id FROM volume WHERE id = "%s"`, volumeID)
	var i europa.Volume
	err := testDB.Get(&i, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &i, nil
}
