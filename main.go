// vim: ts=4 sts=4 sw=4
package main

import (
	"fmt"
	"time"
	"github.com/ofavre/calc-graph/executor"
)

type Data		float32
type ResChan	chan Result
type DataChan	chan Data

type Node interface {
	Out() DataChan
}

type Result struct {
	Node	Node
	value	Data
}


//////////////
// Executor //
//////////////

func LoopNode(tExecutor executor.Executor, node Node) Node {
	runer := node.(executor.Runer)
	tExecutor.Loop(runer)
	// Fluent
	return node
}

////////////////
// SourceNode //
////////////////

var _ Node = (*SourceNode)(nil)
var _ executor.Runer = (*SourceNode)(nil)

type SourceNode struct {
	out	DataChan
	val	Data
}

func NewSourceNode(val Data) *SourceNode {
	return &SourceNode{make(DataChan), val}
}

func (node SourceNode) Run(quitChan executor.QuitChan) {
	select {
		case <-quitChan:
		case node.out <- node.val:
	}
}

func (node SourceNode) Out() DataChan {
	return node.out
}

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


/////////////
// AddNode //
/////////////

var _ Node = (*AddNode)(nil)
var _ executor.Runer = (*AddNode)(nil)

type AddNode struct {
	out,
	in1,
	in2	DataChan
}

func NewAddNode(in1, in2 Node) *AddNode {
	return &AddNode{make(DataChan), in1.Out(), in2.Out()};
}

func (node AddNode) Run(quitChan executor.QuitChan) {
	var sum Data
	var in1, in2, out DataChan
	reset := func() {
		sum = 0
		in1 = node.in1
		in2 = node.in2
		out = nil
	}
	reset()
	CalcLoop: for {
		select {
			case <-quitChan:
				break CalcLoop
			case tmp := <-in1:
				sum += tmp
				in1 = nil
				if in2 == nil {
					out = node.out
				}
			case tmp := <-in2:
				sum += tmp
				in2 = nil
				if in2 == nil {
					out = node.out
				}
			case out <- sum:
				reset()
		}
	}
}

func (node AddNode) Out() DataChan {
	return node.out
}

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

//////////
// Main //
//////////

func main() {
	tExecutor := executor.New()

	obsChan := GoPrintObs(tExecutor)

	var v1 Node = LoopNode(tExecutor, NewSourceNode(1.0))
	fmt.Printf("%#v\n", v1)
	var o1 Node = LoopNode(tExecutor, NewObserverNode(obsChan, v1))
	fmt.Printf("%#v\n", o1)
	var v2 Node = LoopNode(tExecutor, NewSourceNode(2.0))
	fmt.Printf("%#v\n", v2)
	var o2 Node = LoopNode(tExecutor, NewObserverNode(obsChan, v2))
	fmt.Printf("%#v\n", o2)
	var c1 Node = LoopNode(tExecutor, NewAddNode(o1, o2))
	fmt.Printf("%#v\n", c1)
	var o3 Node = LoopNode(tExecutor, NewObserverNode(obsChan, c1))
	fmt.Printf("%#v\n", o3)
	var s1 Node = LoopNode(tExecutor, NewSuckNode(o3))
	fmt.Printf("%#v\n", s1)

	time.Sleep(2 * time.Millisecond)
	tExecutor.Interrupt()
	fmt.Println("Quitting")
}
