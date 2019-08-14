package main

import (
	"flag"
	"fmt"
	"go-projects/customcd/customcd"
	"runtime"
	"time"
)

func main() {
	t := time.Now()
	flag.Parse()
	prefix := flag.Arg(0)
	customcd.SearchPath(prefix)
	fmt.Println(runtime.NumCPU())
	fmt.Println("time taken", time.Since(t))
}
