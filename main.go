package main

import (
	"github.com/eray-can/cycle-detect/runner"
	_ "net/http/pprof"
)

func main() {
	engine := runner.NewEngine()
	engine.Run()
	defer engine.Close()
}
