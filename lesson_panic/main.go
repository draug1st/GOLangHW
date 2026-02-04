package main

import (
	"errors"
	"fmt"
)

var panicError = errors.New("error panic f")

func main() {
	fmt.Println("Hello, World!")
	defer fmt.Println("defer defer 1")
	defer fmt.Println("defer defer 2")
	defer recoverPanic()

	err := safeCall(f)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("end")
}

func recoverPanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if ok && errors.Is(err, panicError) {
			fmt.Println("recovered from panic:", err)
		}
	}

	fmt.Println("close db connection")
}

func f() {
	fmt.Println("f start")
	defer fmt.Println("defer f")
	panic(struct{ Code int }{Code: 100})
	fmt.Println("f end")
}

func safeCall(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error: %v", r)
		}
	}()

	return fn()
}
