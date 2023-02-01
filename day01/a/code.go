package main

import (
	"fmt"
	"os"
)

func main() {
	input, err := os.ReadFile("./input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(input))
}
