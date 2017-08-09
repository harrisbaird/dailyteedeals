package utils

import (
	"io"
	"log"
	"strings"

	"github.com/harrisbaird/dailyteedeals/config"
	minio "github.com/minio/minio-go"
)

func NewMinioConnection() *MinioConnection {
	client, err := minio.New(config.App.AWSS3Endpoint, config.App.AWSAccessKeyID, config.App.AWSSecretAccessKey, config.App.AWSS3Secure)
	if err != nil {
		panic(err)
	}

	conn := &MinioConnection{Client: client, Bucket: config.App.AWSS3Bucket}
	conn.TestConfig()

	return conn
}

func NewMinioTestConnection() *MinioConnection {
	client, err := minio.New("localhost:9000", "access_key", "secret_key", false)
	if err != nil {
		panic(err)
	}

	bucket := RandString(15)
	err = client.MakeBucket(bucket, "us-east-1")
	if err != nil {
		log.Printf("NewMinioTestConnection: MakeBucket returned error: %v", err)
	}

	return &MinioConnection{Client: client, Bucket: bucket, TestMode: true}
}

// RunMinioTest creates a minio test bucket and ensures the bucket is
// removed, even if the test panics.
func RunMinioTest(fn func(*MinioConnection)) {
	conn := NewMinioTestConnection()
	defer conn.Clean()

	fn(conn)
}

type MinioConnection struct {
	Client   *minio.Client
	Bucket   string
	TestMode bool
}

// PutObject creates an object in a bucket.
func (conn *MinioConnection) PutObject(objectName string, reader io.Reader, contentType string) error {
	_, err := conn.Client.PutObject(conn.Bucket, objectName, reader, contentType)
	return err
}

// RemoveObject removes an object from a bucket.
func (conn *MinioConnection) RemoveObject(objectName string) error {
	return conn.Client.RemoveObject(conn.Bucket, objectName)
}

// Clean removes all files from a bucket and deletes the bucket.
func (conn *MinioConnection) Clean() {
	if !conn.TestMode {
		log.Println("Not in test mode, not removing bucket")
		return
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	// Bucket needs to be empty before being removed
	for object := range conn.Client.ListObjects(conn.Bucket, "", true, doneCh) {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		if err := conn.Client.RemoveObject(conn.Bucket, object.Key); err != nil {
			log.Fatalf("Failed to remove %s, error: %v", object.Key, err)
		}
	}

	if err := conn.Client.RemoveBucket(conn.Bucket); err != nil {
		log.Println(err)
	}
}

// TestConfig checks if the minio settings are valid
func (conn *MinioConnection) TestConfig() {
	conn.RemoveObject("connection_test")
	if err := conn.PutObject("connection_test", strings.NewReader(""), "plain/text"); err != nil {
		log.Fatalf("Unable to verify s3 config - %s", err.Error())
	}

	log.Println("S3 config OK")
}
