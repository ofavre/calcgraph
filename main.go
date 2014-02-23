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
	var o1 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, v1))
	fmt.Printf("%v\n", o1)
	var v2 nodes.Node = nodes.LoopNode(executor, nodes.NewConstantNode(uint(2)))
	fmt.Printf("%v\n", v2)
	var o2 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, v2))
	fmt.Printf("%v\n", o2)
	var c1 nodes.Node = nodes.LoopNode(executor, nodes.NewAddNode(nil, o1, o2))
	fmt.Printf("%v\n", c1)
	var o3 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, c1))
	fmt.Printf("%v\n", o3)
	var s1 nodes.Node = nodes.LoopNode(executor, nodes.NewSinkNode(o3))
	fmt.Printf("%v\n", s1)

	time.Sleep(2 * time.Millisecond)
	executor.Interrupt()
	fmt.Println("Quitting")
}
