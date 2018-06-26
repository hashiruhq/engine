package main

import (
	"runtime"
	"trading_engine/cmd"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd.Execute()
}
