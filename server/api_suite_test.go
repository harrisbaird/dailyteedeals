package server_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harrisbaird/dailyteedeals/migrations"
)

func TestMain(m *testing.M) {
	// Silence gin debug output
	gin.SetMode(gin.ReleaseMode)

	migrations.RunTest()
	retCode := m.Run()
	os.Exit(retCode)
}
