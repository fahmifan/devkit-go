package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/devkit/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
