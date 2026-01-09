package main

import (
	"errors"
	"fmt"
)

var CloseDbError = errors.New("close db error")
var CloseRedisError = errors.New("close redis error")

func closeDb() error {
	return CloseDbError
}

func closeRedis() error {
	return CloseRedisError
}

func GracefulShutdown() error {
	var closeErrors []error
	var err error

	err = closeDb()
	if err != nil {
		closeErrors = append(closeErrors, err)
	}
	err = closeRedis()
	if err != nil {
		closeErrors = append(closeErrors, err)
	}

	return errors.Join(closeErrors...)
}

func main() {
	err := GracefulShutdown()
	if err != nil {

		if errors.Is(err, CloseDbError) {
			fmt.Println("Close db error")
		}

		panic(err)
	}
}
