// vim: ts=4 sts=4 sw=4
package main

import (
	"fmt"
	"time"
	"sync"
)

type Data		float32
type ResChan	chan Result
type DataChan	chan Data

type Node interface {
	GetOut() DataChan
}

type Result struct {
	Node	Node
	value	Data
}


//////////////
// Executor //
//////////////

type QuitChan chan bool

type InterruptibleRunable interface {
	Run(QuitChan)
}

type Executor struct {
	quitChan	QuitChan
	waitGroup	sync.WaitGroup
}

func NewExecutor() Executor {
	return Executor{quitChan: make(QuitChan, 100)}
}

func (executor Executor) Run(routine func(QuitChan)) func(QuitChan) {
	executor.waitGroup.Add(1)
	go func() {
		routine(executor.quitChan)
		executor.waitGroup.Done()
	}()
	// Fluent
	return routine
}

func LoopInterruptible(executor Executor, runable InterruptibleRunable) InterruptibleRunable {
	executor.Run(func(quitChan QuitChan) {
		RunLoop: for {
			select {
				case <-quitChan:
					break RunLoop
				default:
					runable.Run(quitChan)
			}
		}
	})
	// Fluent
	return runable
}

func LoopNode(executor Executor, node Node) Node {
	switch runable := node.(type) {
		case InterruptibleRunable:
			LoopInterruptible(executor, runable)
	}
	return node
}

func (executor Executor) WaitInterrupt() {
	// Wait for the WaitGroup to be entirely closed,
	// then close the quitChan
	go func() {
		executor.waitGroup.Wait()
		close(executor.quitChan)
	}()
	defer func() {
		// Expected:
		//   (runtime.errorCString) runtime error: send on closed channel,
		// when the above goroutine has finished waiting for the WaitGroup
		recover()
	}()
	// The following loop breaks in panic when channel closed,
	// and exits the function cleanly
	for {
		executor.quitChan <- true
	}
}


////////////////
// SourceNode //
////////////////

type SourceNode struct {
	out	DataChan
	val	Data
}

func NewSourceNode(val Data) *SourceNode {
	return &SourceNode{make(DataChan), val}
}

func (node SourceNode) Run(quitChan QuitChan) {
	select {
		case <-quitChan:
		case node.out <- node.val:
	}
}

func (node SourceNode) GetOut() DataChan {
	return node.out
}

//////////////
// SuckNode //
//////////////

type SuckNode struct {
	inData	DataChan
}

func NewSuckNode(in Node) *SuckNode {
	return &SuckNode{in.GetOut()}
}

func (node SuckNode) Run(quitChan QuitChan) {
	select {
		case <-quitChan:
		case <-node.inData:
	}
}

func (node SuckNode) GetOut() DataChan {
	return nil
}


/////////////
// AddNode //
/////////////

type AddNode struct {
	out,
	in1,
	in2	DataChan
}

func NewAddNode(in1, in2 Node) *AddNode {
	return &AddNode{make(DataChan), in1.GetOut(), in2.GetOut()};
}

func (node AddNode) Run(quitChan QuitChan) {
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

func (node AddNode) GetOut() DataChan {
	return node.out
}

//////////////
// Observer //
//////////////

type ObserverNode struct {
	out		DataChan
	inData	DataChan
	inNode	Node
	obsChan	ResChan
}

func NewObserverNode(obsChan ResChan, in Node) *ObserverNode {
	return &ObserverNode{make(DataChan), in.GetOut(), in, obsChan}
}

func (node ObserverNode) Run(quitChan QuitChan) {
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

func (node ObserverNode) GetOut() DataChan {
	return node.out
}

func GoPrintObs(executor Executor) ResChan {
	resChan := make(ResChan)
	executor.Run(func(quitChan QuitChan) {
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
	executor := NewExecutor()

	obsChan := GoPrintObs(executor)

	var v1 Node = LoopNode(executor, NewSourceNode(1.0))
	fmt.Printf("%#v\n", v1)
	var o1 Node = LoopNode(executor, NewObserverNode(obsChan, v1))
	fmt.Printf("%#v\n", o1)
	var v2 Node = LoopNode(executor, NewSourceNode(2.0))
	fmt.Printf("%#v\n", v2)
	var o2 Node = LoopNode(executor, NewObserverNode(obsChan, v2))
	fmt.Printf("%#v\n", o2)
	var c1 Node = LoopNode(executor, NewAddNode(o1, o2))
	fmt.Printf("%#v\n", c1)
	var o3 Node = LoopNode(executor, NewObserverNode(obsChan, c1))
	fmt.Printf("%#v\n", o3)
	var s1 Node = LoopNode(executor, NewSuckNode(o3))
	fmt.Printf("%#v\n", s1)

	time.Sleep(2 * time.Millisecond)
	executor.WaitInterrupt()
	fmt.Println("Quitting")
}
