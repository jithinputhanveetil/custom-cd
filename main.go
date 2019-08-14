package main

import (
	"flag"

	"github.com/jithinputhanveetil/custom-cd/customcd"
)

func main() {
	flag.Parse()
	prefix := flag.Arg(0)
	customcd.SearchPath(prefix)
}
