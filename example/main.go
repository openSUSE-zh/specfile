package main

import (
	"fmt"
	"os"

	"github.com/openSUSE-zh/specfile"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()
	if err != nil {
		panic(err)
	}

	parser, err := specfile.NewParser(f)
	if err != nil {
		panic(err)
	}
	parser.Parse()
	fmt.Println(parser.Spec.Dependencies)
}
