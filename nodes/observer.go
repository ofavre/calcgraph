package nodes

import (
	"fmt"
	"github.com/ofavre/calcgraph/executor"
)

//////////////
// Observer //
//////////////

var _ Node = (*ObserverNode)(nil)
var _ executor.Runer = (*ObserverNode)(nil)

type ObserverNode struct {
	out		DataChan
	inData	DataChan
	inNode	Node
	obsChan	ResChan
}

func NewObserverNode(obsChan ResChan, in Node) *ObserverNode {
	return &ObserverNode{make(DataChan), in.Out(), in, obsChan}
}

func (node ObserverNode) Run(quitChan executor.QuitChan) {
	var data Data
	var ok bool
	select {
		case <-quitChan:
			close(node.out)
			return
		case data, ok = <-node.inData:
			if !ok {
				return
			}
	}
	select {
		case <-quitChan:
			close(node.out)
			return
		case node.obsChan <- Result{node.inNode, data}:
	}
	select {
		case <-quitChan:
			close(node.out)
			return
		case node.out <- data:
	}
}

func (node ObserverNode) Out() DataChan {
	return node.out
}

func (node ObserverNode) String() string {
	return fmt.Sprintf("%T{%v}", node, node.inNode)
}

func GoPrintObs(tExecutor *executor.Executor) ResChan {
	resChan := make(ResChan)
	tExecutor.Run(func(quitChan executor.QuitChan) {
		fmt.Println("Observer started")
		ReadLoop: for {
			select {
				case <-quitChan:
					break ReadLoop
				case res, ok := <-resChan:
					if !ok {
						break ReadLoop
					}
					fmt.Printf("Obs: %v -> %T(%#v)\n", res.Node, res.Value, res.Value)
			}
		}
		fmt.Println("Observer stopped")
	})
	return resChan
}
