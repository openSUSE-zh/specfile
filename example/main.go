package main

import (
	"fmt"
	"os"

	specfile "github.com/openSUSE/specfile"
)

func main() {
	f, err := os.Open("../test/ffmpeg-4.spec")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines := specfile.ReadLines(f)
	nodes := specfile.NewLexer(lines)
	fmt.Println(nodes)
}
