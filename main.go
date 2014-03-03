// vim: ts=4 sts=4 sw=4
package main

import (
	"fmt"
	"time"
	"github.com/ofavre/calcgraph/executor"
	"github.com/ofavre/calcgraph/nodes"
)

//////////
// Main //
//////////

func main() {
	executor := executor.New()

	obsChan := nodes.GoPrintObs(executor)

	var v1 nodes.Node = nodes.LoopNode(executor, nodes.NewConstantNode(int(1)))
	fmt.Printf("%v\n", v1)
	var v2 nodes.Node = nodes.LoopNode(executor, nodes.NewConstantNode(int(2)))
	fmt.Printf("%v\n", v2)
	var c1 nodes.Node = nodes.LoopNode(executor, nodes.NewAddNode(nil, v1, v2))
	fmt.Printf("%v\n", c1)
	var c2 nodes.Node = nodes.LoopNode(executor, nodes.NewSubNode(nil, v1, c1))
	fmt.Printf("%v\n", c2)
	var c3 nodes.Node = nodes.LoopNode(executor, nodes.NewMulNode(nil, c1, c2))
	fmt.Printf("%v\n", c2)
	var c4 nodes.Node = nodes.LoopNode(executor, nodes.NewFanInNode(c1, c2, c3))
	fmt.Printf("%v\n", c4)
	var o nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, c4))
	fmt.Printf("%v\n", o)
	var s1 nodes.Node = nodes.LoopNode(executor, nodes.NewSinkNode(o))
	fmt.Printf("%v\n", s1)

	time.Sleep(2 * time.Millisecond)
	executor.Interrupt()
	fmt.Println("Quitting")
}
