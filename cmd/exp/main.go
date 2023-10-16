package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()
	if errors.Is(err, ErrNotFound) {
		fmt.Println("err is not found")
	}

	if errors.As(err, &ErrNotFound) {
		fmt.Println("err is not found again")
		return
	}
}

var ErrNotFound = errors.New("not found")

func A() error {
	return ErrNotFound
}

func B() error {
	err := A()
	return fmt.Errorf("b: %w", err)
}
