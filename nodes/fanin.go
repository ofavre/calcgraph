package nodes

import (
	"fmt"
	"reflect"
	"github.com/ofavre/calcgraph/executor"
)

///////////////
// FanInNode //
///////////////

var _ Node = (*FanInNode)(nil)
var _ executor.Runer = (*FanInNode)(nil)

type FanInNode struct {
	out				DataChan
	inNodes			[]Node
	selectCases		[]reflect.SelectCase
	remaining		int
	lastQuitChan	executor.QuitChan
}

func NewFanInNode(inNodes ...Node) *FanInNode {
	selectCases := make([]reflect.SelectCase, 1+len(inNodes))
	var nilValue reflect.Value
	for i, node := range inNodes {
		selectCases[1+i] = reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(node.Out()), nilValue}
	}
	return &FanInNode{make(DataChan), inNodes, selectCases, len(inNodes), nil};
}

func (node *FanInNode) Run(quitChan executor.QuitChan) {
	if node.remaining <= 0 {
		close(node.out)
		<-quitChan
		return
	}
	if (node.lastQuitChan != quitChan) {
		var nilValue reflect.Value
		node.selectCases[0] = reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(quitChan), nilValue}
	}
	chosen, val, ok := reflect.Select(node.selectCases)
	if chosen == 0 { // quitChan
		close(node.out)
	} else if ok {
		select {
			case <-quitChan:
				close(node.out)
			case node.out <- Data(val.Interface()):
		}
	} else {
		node.selectCases[chosen].Chan = reflect.ValueOf(nil)
	}
}

func (node FanInNode) Out() DataChan {
	return node.out
}

func (node FanInNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNodes)
}
