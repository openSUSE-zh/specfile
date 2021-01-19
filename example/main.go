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
	err = parser.Parse()
	if err != nil {
		panic(err)
	}

	for _, v := range parser.Spec.Subpackages {
		fmt.Println(v)
	}
}
