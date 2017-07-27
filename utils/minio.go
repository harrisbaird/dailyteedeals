package utils

import (
	"log"

	"github.com/harrisbaird/dailyteedeals/config"
	minio "github.com/minio/minio-go"
)

func NewMinioConnection() *MinioConnection {
	client, err := minio.New(config.App.AWSS3Endpoint, config.App.AWSAccessKeyID, config.App.AWSSecretAccessKey, config.App.AWSS3Secure)
	if err != nil {
		panic(err)
	}
	return &MinioConnection{Client: client, Bucket: config.App.AWSS3Bucket}
}

func NewMinioTestConnection() *MinioConnection {
	client, err := minio.New("localhost:9000", "access_key", "secret_key", false)
	if err != nil {
		panic(err)
	}

	bucket := RandString(15)
	err = client.MakeBucket(bucket, "us-east-1")
	if err != nil {
		log.Printf("NewMinioTestConnection: MakeBucket returned error: %v\n", err)
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
			log.Fatalf("Failed to remove %s, error: %v\n", object.Key, err)
		}
	}

	if err := conn.Client.RemoveBucket(conn.Bucket); err != nil {
		log.Println(err)
	}
}
