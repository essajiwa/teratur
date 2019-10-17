package boot

import "fmt"

// GRPC will load configuration, do DI and then start the gRPC server
func GRPC() error {
	fmt.Println("gRPC server started")
	return nil
}
