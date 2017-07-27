package models_test

import (
	"os"
	"testing"

	"github.com/harrisbaird/dailyteedeals/migrations"
	"github.com/harrisbaird/dailyteedeals/utils"
)

func init() {
	utils.SetHTTPTestMode()
}

func TestMain(m *testing.M) {
	migrations.RunTest()
	retCode := m.Run()
	os.Exit(retCode)
}
