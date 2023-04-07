package main

import (
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/dmcli/cmd"
)

func main() {

	definition.Load()
	cmd.Execute()
}
