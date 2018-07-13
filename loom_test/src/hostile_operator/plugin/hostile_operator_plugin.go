package main

import (
	"hostile_operator"

	"github.com/loomnetwork/go-loom/plugin"
)

var Contract = hostile_operator.Contract

func main() {
	plugin.Serve(Contract)
}
