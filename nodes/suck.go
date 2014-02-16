package nodes

import (
	"github.com/ofavre/calcgraph/executor"
)

//////////////
// SuckNode //
//////////////

var _ Node = (*SuckNode)(nil)
var _ executor.Runer = (*SuckNode)(nil)

type SuckNode struct {
	inData	DataChan
}

func NewSuckNode(in Node) *SuckNode {
	return &SuckNode{in.Out()}
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
