package main

import (
	"errors"
	"fmt"
	"strconv"
)

func ParsePositiveInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if val <= 0 {
		return 0, errors.New("must be positive")
	}
	return val, nil
}

func SumTwo(a, b string) (int, error) {
	valA, err := ParsePositiveInt(a)
	if err != nil {
		return 0, err
	}

	valB, err := ParsePositiveInt(b)
	if err != nil {
		return 0, err
	}

	return valA + valB, nil
}

func main() {
	res, err := SumTwo("10", "20")
	fmt.Printf("SumTwo(\"10\",\"20\") -> res: %d, err: %v\n", res, err)

	res, err = SumTwo("x", "20")
	fmt.Printf("SumTwo(\"x\",\"20\") -> res: %d, err: %v\n", res, err)

	res, err = SumTwo("-1", "20")
	fmt.Printf("SumTwo(\"-1\",\"20\") -> res: %d, err: %v\n", res, err)
}
