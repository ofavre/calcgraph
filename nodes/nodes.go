package nodes

import (
	"github.com/ofavre/calcgraph/executor"
)

type Data		interface{}
type ResChan	chan Result
type DataChan	chan Data

type Node interface {
	executor.Runer
	Out() DataChan
}

type Result struct {
	Node	Node
	Value	Data
}


//////////////
// Executor //
//////////////

func LoopNode(tExecutor *executor.Executor, node Node) Node {
	tExecutor.Loop(node)
	// Fluent
	return node
}
