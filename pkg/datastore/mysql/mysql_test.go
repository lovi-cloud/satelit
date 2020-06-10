package mysql_test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/internal/testutils"
	"github.com/whywaita/satelit/pkg/europa"
)

const (
	testUUID = "90dd6cd4-b3e4-47f3-9af5-47f78efc8fc7"
)

var testImage = europa.BaseImage{
	UUID:          uuid.FromStringOrNil(testUUID),
	Name:          "test-image",
	Description:   "test-image-description",
	CacheVolumeID: "test-volume-id",
}

func TestMain(m *testing.M) {
	os.Exit(testutils.IntegrationTestRunner(m))
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

	images, err := testDatastore.GetImages()
	if err != nil {
		t.Fatalf("failed to get image: %s", err)
	}
	if len(images) != 1 {
		t.Fatalf("unexpected images value, image count: (expected: 1, actual: %d)", len(images))
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

	err = testDatastore.DeleteImage(testUUID)
	if err != nil {
		t.Fatalf("failed to delete image: %s", err)
	}

	_, err = getImageFromSQL(testDB, testUUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("failed to get image from sql: %s", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("unexpected images value, this test is after DeleteImage but there are image")
	}
}

func getImageFromSQL(testDB *sqlx.DB, imageID string) (*europa.BaseImage, error) {
	query := fmt.Sprintf(`SELECT * FROM image WHERE BIN_TO_UUID(uuid) = "%s"`, imageID)
	var i europa.BaseImage
	err := testDB.Get(&i, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s", err)
	}

	return &i, nil
}