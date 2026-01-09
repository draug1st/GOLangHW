package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(uuid.New().String())
}
