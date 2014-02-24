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
	var v2 nodes.Node = nodes.LoopNode(executor, nodes.NewConstantNode(int(2)))
	fmt.Printf("%v\n", v2)
	var o2 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, v2))
	fmt.Printf("%v\n", o2)
	var c1 nodes.Node = nodes.LoopNode(executor, nodes.NewAddNode(nil, o1, o2))
	fmt.Printf("%v\n", c1)
	var o3 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, c1))
	fmt.Printf("%v\n", o3)
	var c2 nodes.Node = nodes.LoopNode(executor, nodes.NewSubNode(nil, o1, o3))
	fmt.Printf("%v\n", c2)
	var o4 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, c2))
	fmt.Printf("%v\n", o4)
	var c3 nodes.Node = nodes.LoopNode(executor, nodes.NewMulNode(nil, o3, o4))
	fmt.Printf("%v\n", c2)
	var o5 nodes.Node = nodes.LoopNode(executor, nodes.NewObserverNode(obsChan, c3))
	fmt.Printf("%v\n", o4)
	var s1 nodes.Node = nodes.LoopNode(executor, nodes.NewSinkNode(o5))
	fmt.Printf("%v\n", s1)

	time.Sleep(2 * time.Millisecond)
	executor.Interrupt()
	fmt.Println("Quitting")
}
