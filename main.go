package main // import "gitlab.com/around25/products/matching-engine"

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
