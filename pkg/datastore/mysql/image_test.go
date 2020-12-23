package mysql_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/satelit/internal/testutils"
	"github.com/lovi-cloud/satelit/pkg/europa"
)

var testImage = europa.BaseImage{
	UUID:          uuid.FromStringOrNil(testUUID),
	Name:          "test-image",
	Description:   "test-image-description",
	CacheVolumeID: "test-volume-id",
}

func TestMySQL_PutImage(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %s", err)
	}

	i, err := getImageFromSQL(testDB, testUUID)
	if err != nil {
		t.Fatalf("failed to get image from sql: %s", err)
	}

	ok, values := testutils.CompareStruct(testImage, i)
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}
}

func TestMySQL_GetImage(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %s", err)
	}

	image, err := testDatastore.GetImage(testImage.UUID)
	if err != nil {
		t.Fatalf("failed to get image: %s", err)
	}
	ok, values := testutils.CompareStruct(testImage, image)
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}
}

func TestMySQL_ListImage(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %+v", err)
	}

	images, err := testDatastore.ListImage()
	if err != nil {
		t.Fatalf("failed to list image: %+v", err)
	}
	if len(images) != 1 {
		t.Fatalf("unexpected images value, image count: (expected: 1, actual: %d)", len(images))
	}
	ok, values := testutils.CompareStruct(testImage, images[0])
	if !ok {
		t.Fatalf("unexpected values, field name: %s, input: %s, output: %s", values[0], values[1], values[2])
	}
}

func TestMySQL_DeleteImage(t *testing.T) {
	testDatastore, teardown := testutils.GetTestDatastore()
	defer teardown()
	testDB, _ := testutils.GetTestDB()

	err := testDatastore.PutImage(testImage)
	if err != nil {
		t.Fatalf("failed to put image: %s", err)
	}

	err = testDatastore.DeleteImage(uuid.FromStringOrNil(testUUID))
	if err != nil {
		t.Fatalf("failed to delete image: %s", err)
	}

	b, err := getImageFromSQL(testDB, testUUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("failed to get image from sql: %s", err)
	}
	if b != nil {
		t.Fatalf("unexpected images value, this test is after DeleteImage but there are image")
	}

	// sql.ErrNoRows is correct
}

func getImageFromSQL(testDB *sqlx.DB, imageID string) (*europa.BaseImage, error) {
	var i europa.BaseImage
	query := `SELECT uuid, name, description, volume_id, created_at, updated_at FROM image WHERE uuid = ?`
	stmt, err := testDB.Preparex(query)
	err = stmt.Get(&i, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &i, nil
}
