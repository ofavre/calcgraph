package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

//////////////
// SuckNode //
//////////////

var _ Node = (*SuckNode)(nil)
var _ executor.Runer = (*SuckNode)(nil)

type SuckNode struct {
	inNode	Node
	inData	DataChan
}

func NewSuckNode(in Node) *SuckNode {
	return &SuckNode{in, in.Out()}
}

func (node SuckNode) Run(quitChan executor.QuitChan) {
	select {
		case <-quitChan:
		case <-node.inData:
	}
}

func (node SuckNode) Out() DataChan {
	return nil
}

func (node SuckNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNode)
}
