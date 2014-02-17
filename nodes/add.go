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
	out			DataChan
	inNodes		[]Node
	assembler	*Assembler
}

func NewAddNode(inNodes ...Node) *AddNode {
	return &AddNode{make(DataChan), inNodes, NewAssembler(inNodes...)};
}

func (node AddNode) Run(quitChan executor.QuitChan) {
	go node.assembler.Run(quitChan)
	assemblerOut := node.assembler.Out()
	var sum Data
	var out DataChan = nil
	CalcLoop: for {
		select {
			case <-quitChan:
				break CalcLoop
			case inputs := <-assemblerOut:
				out = node.out
				sum = 0
				for _, val := range inputs {
					sum += val
				}
			case out <- sum:
				break CalcLoop
		}
	}
}

func (node AddNode) Out() DataChan {
	return node.out
}

func (node AddNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNodes)
}
