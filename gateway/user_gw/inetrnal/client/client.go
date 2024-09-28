package client

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcClientServerConnection(server_address string) *grpc.ClientConn {
	conn, err := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func MinioClient(minioServerAddress, minioAccessKey, minioSecretKey string) *minio.Client {
	client, err := minio.New(minioServerAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("could not connect to the minio server: %v", err)
	}
	return client
}
