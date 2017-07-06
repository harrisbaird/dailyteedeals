package utils

import (
	"go/build"
	"path"
)

func ProjectRootPath(elem ...string) string {
	rootPath := []string{build.Default.GOPATH, "src", "github.com", "harrisbaird", "dailyteedeals"}
	rootPath = append(rootPath, elem...)
	return path.Join(rootPath...)
}
