package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

//////////////////
// ConstantNode //
//////////////////

var _ Node = (*ConstantNode)(nil)
var _ executor.Runer = (*ConstantNode)(nil)

type ConstantNode struct {
	out	DataChan
	val	Data
}

func NewConstantNode(val Data) *ConstantNode {
	return &ConstantNode{make(DataChan), val}
}

func (node ConstantNode) Run(quitChan executor.QuitChan) {
	select {
		case <-quitChan:
		case node.out <- node.val:
	}
}

func (node ConstantNode) Out() DataChan {
	return node.out
}

func (node ConstantNode) String() string {
	return fmt.Sprintf("%T{%T(%v)}", node, node.val, node.val)
}
