package main

import (
	"fmt"
	"os"
)

func main() {
	if err := NewCmd().Exec(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
