package main

import (
	"fmt"

	"github.com/essajiwa/teratur/internal/boot"
)

func main() {
	if err := boot.HTTP(); err != nil {
		fmt.Println("[HTTP] failed to boot http server due to " + err.Error())
	}
}
