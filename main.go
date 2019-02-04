package main

import (
	"runtime"

	"gitlab.com/around25/products/matching-engine/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd.Execute()
}
