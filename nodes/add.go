package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

/////////////
// AddNode //
/////////////

var _ Node = (*AddNode)(nil)
var _ executor.Runer = (*AddNode)(nil)

type AddNode struct {
	out DataChan
	inNode1,
	inNode2 Node
	in1,
	in2	DataChan
}

func NewAddNode(in1, in2 Node) *AddNode {
	return &AddNode{make(DataChan), in1, in2, in1.Out(), in2.Out()};
}

func (node AddNode) Run(quitChan executor.QuitChan) {
	var sum Data
	var in1, in2, out DataChan
	reset := func() {
		sum = 0
		in1 = node.in1
		in2 = node.in2
		out = nil
	}
	reset()
	CalcLoop: for {
		select {
			case <-quitChan:
				break CalcLoop
			case tmp := <-in1:
				sum += tmp
				in1 = nil
				if in2 == nil {
					out = node.out
				}
			case tmp := <-in2:
				sum += tmp
				in2 = nil
				if in2 == nil {
					out = node.out
				}
			case out <- sum:
				reset()
		}
	}
}

func (node AddNode) Out() DataChan {
	return node.out
}

func (node AddNode) String() string {
	return fmt.Sprintf("%T{%v + %v}", node, node.inNode1, node.inNode2)
}
