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
	var data Data;
	select {
		case <-quitChan:
			return
		case data = <-node.inData:
	}
	select {
		case <-quitChan:
		case node.obsChan <- Result{node.inNode, data}:
	}
	select {
		case <-quitChan:
		case node.out <- data:
	}
}

func (node ObserverNode) Out() DataChan {
	return node.out
}

func GoPrintObs(tExecutor executor.Executor) ResChan {
	resChan := make(ResChan)
	tExecutor.Run(func(quitChan executor.QuitChan) {
		fmt.Println("Observer started")
		ReadLoop: for {
			select {
				case <-quitChan:
					break ReadLoop
				case res := <-resChan:
					fmt.Printf("Obs: %#v\n", res)
			}
		}
		// FIXME Not always visible in the console, but always visible in a pipe…
		fmt.Println("Observer stopped")
	})
	return resChan
}

