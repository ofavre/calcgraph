package nodes

import (
	"fmt"
	"reflect"
	"github.com/ofavre/calcgraph/executor"
	"github.com/ofavre/calcgraph/numeric"
)

/////////////
// SubNode //
/////////////

var _ Node = (*SubNode)(nil)
var _ executor.Runer = (*SubNode)(nil)

type SubNode struct {
	out				DataChan
	inNodes			[]Node
	typeEnforced	reflect.Type
	assembler		*Assembler
}

func NewSubNode(typeEnforced reflect.Type, inNodes ...Node) *SubNode {
	return &SubNode{make(DataChan), inNodes, typeEnforced, NewAssembler(typeEnforced, inNodes...)};
}

func (node SubNode) Run(quitChan executor.QuitChan) {
	go node.assembler.Run(quitChan)
	assemblerOut := node.assembler.Out()
	var diff Data
	var out DataChan = nil
	CalcLoop: for {
		select {
			case <-quitChan:
				break CalcLoop
			case inputs := <-assemblerOut:
				out = node.out
				if len(inputs) == 0 {
					if node.typeEnforced != nil {
						diff = numeric.ZeroOfType(node.typeEnforced)
					} else {
						diff = 0
					}
				} else {
					diff = inputs[0]
					for _, val := range inputs[1:] {
						diff = numeric.Sub(diff, val)
					}
				}
			case out <- diff:
				break CalcLoop
		}
	}
}

func (node SubNode) Out() DataChan {
	return node.out
}

func (node SubNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNodes)
}
