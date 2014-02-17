package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

//////////////
// SinkNode //
//////////////

var _ Node = (*SinkNode)(nil)
var _ executor.Runer = (*SinkNode)(nil)

type SinkNode struct {
	inNode	Node
	inData	DataChan
}

func NewSinkNode(in Node) *SinkNode {
	return &SinkNode{in, in.Out()}
}

func (node SinkNode) Run(quitChan executor.QuitChan) {
	select {
		case <-quitChan:
		case <-node.inData:
	}
}

func (node SinkNode) Out() DataChan {
	return nil
}

func (node SinkNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNode)
}
