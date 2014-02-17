package nodes

import (
	"github.com/ofavre/calcgraph/executor"
)

type AssembledData		[]Data
type AssembledDataChan	chan AssembledData

type collectChan	chan Data

type Assembler struct {
	inNodes			[]Node
	assembledChan	AssembledDataChan
	collectChan		collectChan
}

func NewAssembler(inNodes ...Node) *Assembler {
	return &Assembler{inNodes, make(AssembledDataChan), make(collectChan, len(inNodes))}
}

func (assembler Assembler) Out() AssembledDataChan {
	return assembler.assembledChan
}

func (assembler Assembler) Run(quitChan executor.QuitChan) {
	missing	:= len(assembler.inNodes)
	results	:= make(AssembledData, missing)
	// Start one worker per input node
	for _, node := range assembler.inNodes {
		go assemblerWorker(quitChan, node, assembler.collectChan)
	}
	// Collect inputs
	MainWorkLoop: for {
		select {
			case <-quitChan:
				break MainWorkLoop
			case val := <-assembler.collectChan:
				missing--
				results[missing] = val
				if missing == 0 {
					break MainWorkLoop
				}
		}
	}
	assembler.assembledChan <- results
}

func assemblerWorker(quitChan executor.QuitChan, node Node, collectChan collectChan) {
	select {
		case <-quitChan:
			return
		// Read value to transmit
		case val := <-node.Out():
			select {
				case <-quitChan:
					return
				// Transmit the value
				case collectChan <- val:
			}
	}
}
