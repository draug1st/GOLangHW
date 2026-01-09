package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNegativeAge = errors.New("negative age")
	ErrTooBigAge   = errors.New("too big age")
	ErrTooYoungAge = errors.New("too young age")
)

type ValidationError struct {
	Field  string
	Reason string
	Err    error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ERROR -> field %s: %s", e.Field, e.Reason)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{
			Field:  "email",
			Reason: "empty email",
		}
	}
	if !strings.Contains(email, "@") {
		return &ValidationError{
			Field:  "email",
			Reason: "no @ in email",
		}
	}
	if !strings.Contains(email, ".") {
		return &ValidationError{
			Field:  "email",
			Reason: "invalid email",
		}
	}
	return nil
}

func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Parse age error: %v", err)
	}
	if n < 0 {
		return 0, ErrNegativeAge
	}
	if n < 18 {
		return 0, ErrTooYoungAge
	}
	if n > 150 {
		return 0, ErrTooBigAge
	}
	return n, nil
}

func main() {

	email := "dasd@a"
	err := ValidateEmail(email)
	if err != nil {
		var validationError *ValidationError
		if errors.As(err, &validationError) {
			fmt.Println(validationError)
			return
		}
		fmt.Println(err)
		return
	}

	curr_age := "asd"
	age, err := ParseAge(curr_age)
	if err != nil {
		var validationError *ValidationError
		if errors.As(err, &validationError) {
			fmt.Println(validationError)
			return
		}
		if errors.Is(err, ErrNegativeAge) {
			fmt.Printf("negative age %s\n", curr_age)
			return
		}
		if errors.Is(err, ErrTooBigAge) {
			fmt.Printf("too big age %s \n", curr_age)
			return
		}

		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			fmt.Println(numErr.Num)
		}

		fmt.Println(err)
		return
	}

	fmt.Println(age)
}
