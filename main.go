package main // import "github.com/hashiruhq/engine"

import (
	"runtime"

	"github.com/hashiruhq/engine/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd.Execute()
}
