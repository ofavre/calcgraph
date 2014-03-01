// vim: ts=4 sts=4 sw=4
package executor

import (
	"sync"
)

//////////////
// Executor //
//////////////

type QuitChan chan bool

type Runer interface {
	Run(QuitChan)
}

type Executor struct {
	quitChan	QuitChan
	waitGroup	*sync.WaitGroup
}

func New() *Executor {
	return &Executor{make(QuitChan), new(sync.WaitGroup)}
}

func (executor *Executor) Run(routine func(QuitChan)) func(QuitChan) {
	executor.waitGroup.Add(1)
	go func() {
		routine(executor.quitChan)
		executor.waitGroup.Done()
	}()
	// Fluent
	return routine
}

func (executor *Executor) Loop(runable Runer) Runer {
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

func (executor *Executor) Wait() {
	// Wait for the WaitGroup to be entirely closed
	executor.waitGroup.Wait()
}

func (executor *Executor) Interrupt() {
	// Closing the QuitChan will unlock every receive operations on it
	close(executor.quitChan)
	// Create another one for future reuse of the Executor
	executor.quitChan = make(QuitChan)
	// Wait for the WaitGroup to be entirely closed,
	executor.Wait()
}
