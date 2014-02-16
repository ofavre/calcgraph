package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

////////////////
// SourceNode //
////////////////

var _ Node = (*SourceNode)(nil)
var _ executor.Runer = (*SourceNode)(nil)

type SourceNode struct {
	out	DataChan
	val	Data
}

func NewSourceNode(val Data) *SourceNode {
	return &SourceNode{make(DataChan), val}
}

func (node SourceNode) Run(quitChan executor.QuitChan) {
	select {
		case <-quitChan:
		case node.out <- node.val:
	}
}

func (node SourceNode) Out() DataChan {
	return node.out
}

func (node SourceNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.val)
}
