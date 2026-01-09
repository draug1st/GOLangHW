package main

import (
	"fmt"

	"github.com/draug1st/GOLangHW/tree/main/greeting"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(uuid.New().String())
	fmt.Println(greeting.Hello())
	fmt.Println(greeting.Bye())
}
