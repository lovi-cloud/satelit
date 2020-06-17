package mysql_test

import (
	"os"
	"testing"

	"github.com/whywaita/satelit/internal/testutils"
)

const (
	testUUID = "90dd6cd4-b3e4-47f3-9af5-47f78efc8fc7"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.IntegrationTestRunner(m))
}
