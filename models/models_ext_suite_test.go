package models_test

import (
	"os"
	"testing"

	"github.com/harrisbaird/dailyteedeals/utils"
)

func init() {
	utils.SetHTTPTestMode()
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}
