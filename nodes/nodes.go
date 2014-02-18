package nodes

import (
	"github.com/ofavre/calcgraph/executor"
)

type Data		float32
type ResChan	chan Result
type DataChan	chan Data

type Node interface {
	Out() DataChan
}

type Result struct {
	Node	Node
	value	Data
}


//////////////
// Executor //
//////////////

func LoopNode(tExecutor *executor.Executor, node Node) Node {
	runer := node.(executor.Runer)
	tExecutor.Loop(runer)
	// Fluent
	return node
}
