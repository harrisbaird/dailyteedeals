package backend_test

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/harrisbaird/dailyteedeals/utils"
)

func init() {
	utils.SetHTTPTestMode()
}

func filesIdentical(file1, file2 *os.File) bool {
	b1, err := ioutil.ReadAll(file1)
	if err != nil {
		panic(err)
	}

	b2, err := ioutil.ReadAll(file2)
	if err != nil {
		panic(err)
	}

	return bytes.Compare(b1, b2) == 0
}
