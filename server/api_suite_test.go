package server_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Silence gin debug output
	gin.SetMode(gin.ReleaseMode)
	retCode := m.Run()
	os.Exit(retCode)
}
