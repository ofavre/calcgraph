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
	tExecutor := executor.New()

	obsChan := nodes.GoPrintObs(tExecutor)

	var v1 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewSourceNode(1.0))
	fmt.Printf("%v\n", v1)
	var o1 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewObserverNode(obsChan, v1))
	fmt.Printf("%v\n", o1)
	var v2 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewSourceNode(2.0))
	fmt.Printf("%v\n", v2)
	var o2 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewObserverNode(obsChan, v2))
	fmt.Printf("%v\n", o2)
	var c1 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewAddNode(o1, o2))
	fmt.Printf("%v\n", c1)
	var o3 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewObserverNode(obsChan, c1))
	fmt.Printf("%v\n", o3)
	var s1 nodes.Node = nodes.LoopNode(tExecutor, nodes.NewSuckNode(o3))
	fmt.Printf("%v\n", s1)

	time.Sleep(2 * time.Millisecond)
	tExecutor.Interrupt()
	fmt.Println("Quitting")
}
