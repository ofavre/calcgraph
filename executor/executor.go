// vim: ts=4 sts=4 sw=4
package executor

import (
	"sync"
)

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
