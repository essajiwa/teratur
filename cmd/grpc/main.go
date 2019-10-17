package main

import (
	"fmt"

	"github.com/essajiwa/teratur/internal/boot"
)

func main() {
	if err := boot.GRPC(); err != nil {
		fmt.Println("[gRPC] failed to boot grpc server due to " + err.Error())
	}
}
