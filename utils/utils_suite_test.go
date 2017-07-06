package utils_test

import (
	"image"

	"github.com/disintegration/imaging"
	. "github.com/harrisbaird/dailyteedeals/utils"
)

func init() {
	SetHTTPTestMode()
}

func LoadImageFixture(elem ...string) image.Image {
	fixturePath := []string{"utils", "testdata"}
	fixturePath = append(fixturePath, elem...)
	file, err := imaging.Open(ProjectRootPath(fixturePath...))
	if err != nil {
		panic(err)
	}
	return file
}
