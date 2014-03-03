package nodes

import (
	"fmt"
	"reflect"
	"github.com/ofavre/calcgraph/executor"
	"github.com/ofavre/calcgraph/numeric"
)

/////////////
// MulNode //
/////////////

var _ Node = (*MulNode)(nil)
var _ executor.Runer = (*MulNode)(nil)

type MulNode struct {
	out				DataChan
	inNodes			[]Node
	typeEnforced	reflect.Type
	assembler		*AssemblerNode
}

func NewMulNode(typeEnforced reflect.Type, inNodes ...Node) *MulNode {
	return &MulNode{make(DataChan), inNodes, typeEnforced, NewAssemblerNode(typeEnforced, inNodes...)};
}

func (node MulNode) Run(quitChan executor.QuitChan) {
	go node.assembler.Run(quitChan)
	assemblerOut := node.assembler.Out()
	var product Data
	var out DataChan = nil
	CalcLoop: for {
		select {
			case <-quitChan:
				close(node.out)
				break CalcLoop
			case val, ok := <-assemblerOut:
				if !ok {
					break CalcLoop
				}
				inputs := val.([]Data)
				out = node.out
				if len(inputs) == 0 {
					if node.typeEnforced != nil {
						product = numeric.ConvertToType(1, node.typeEnforced)
					} else {
						product = 1
					}
				} else {
					product = inputs[0]
					for _, val := range inputs[1:] {
						product = numeric.Mul(product, val)
					}
				}
			case out <- product:
				break CalcLoop
		}
	}
}

func (node MulNode) Out() DataChan {
	return node.out
}

func (node MulNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNodes)
}
