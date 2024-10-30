package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Magowtham/go_file_streaming_server/proto/filestream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConnection, err := grpc.NewClient("127.0.0.1:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Printf("error occurred -> %v\n", err.Error())
		return
	}

	defer grpcConnection.Close()

	client := filestream.NewFileStreamServiceClient(grpcConnection)

	fileRequest := &filestream.FileRequest{
		Filename: "departments.txt",
	}

	stream, err := client.DownloadFile(context.Background(), fileRequest)

	if err != nil {
		fmt.Printf("error occurred -> %v\n", err.Error())
		return
	}

	file, err := os.Create(fileRequest.Filename)

	if err != nil {
		fmt.Printf("error occurred -> %v", err.Error())
		return
	}

	defer file.Close()

	for {

		chunk, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("file downloaded successfully 😎..\n")
			break
		}

		if err != nil {
			fmt.Printf("error occurred -> %v", err.Error())
			return
		}

		var buffer []byte

		for _, value := range chunk.Chunk {
			if value == 0 {
				break
			}

			buffer = append(buffer, value)
		}
		size, err := file.Write(buffer)

		if err != nil {
			fmt.Printf("error occurred -> %v", err.Error())
			return
		}

		fmt.Printf("💀 successfully written %v bytes\n", size)
	}

}
