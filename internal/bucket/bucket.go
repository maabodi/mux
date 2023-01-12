package bucket

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	ctx := context.Background()
	// endpoint := "139.180.221.53:9000"
	// accessKeyID := "minioadmin"
	// secretAccessKey := "minioadmin"
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// minioClient is now setup
	bucket := "images"

	err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucket)
	}

	return minioClient
}

func UploadFile(fileName, filePath, bucketName, fileExt string) (url string, err error) {
	ctx := context.Background()

	minioClient := InitMinio()

	//up
	info, err := minioClient.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType: fileExt,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Print(info)
	log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)

	fileUrl := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, fileName)
	return fileUrl, nil

}

func RemoveFile(fileName string, fileType string) error {
	log.Print("deleting ", fileName)
	ctx := context.Background()
	minioClient := InitMinio()

	err := minioClient.RemoveObject(ctx, fileType, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	log.Printf("Successfully deleted %s", fileName)
	return nil
}
