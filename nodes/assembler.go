package nodes

import (
	"fmt"
	"reflect"
	"github.com/ofavre/calcgraph/executor"
)



type TypeMismatchError struct {
	expectedType	reflect.Type
	actualType		reflect.Type
}

func typeMismatchErrorForTypeOf(expectedType reflect.Type, val interface{}) {
	typeMismatchError(expectedType, reflect.TypeOf(val))
}

func typeMismatchError(expectedType, actualType reflect.Type) {
	panic(TypeMismatchError{expectedType, actualType})
}

func (err TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: %v does not match expected %v", err.actualType, err.expectedType)
}



type AssembledData		[]Data
type AssembledDataChan	chan AssembledData

type PositionedData struct {
	position	int
	data		Data
}
type collectChan	chan PositionedData

type Assembler struct {
	typeEnforced	reflect.Type
	inNodes			[]Node
	assembledChan	AssembledDataChan
	collectChan		collectChan
}

func NewAssembler(typeEnforced reflect.Type, inNodes ...Node) *Assembler {
	return &Assembler{typeEnforced, inNodes, make(AssembledDataChan), make(collectChan, len(inNodes))}
}

func (assembler Assembler) Out() AssembledDataChan {
	return assembler.assembledChan
}

func (assembler Assembler) Run(quitChan executor.QuitChan) {
	missing	:= len(assembler.inNodes)
	results	:= make(AssembledData, missing)
	// Start one worker per input node
	for i, node := range assembler.inNodes {
		go assemblerWorker(quitChan, node, i, assembler.collectChan, assembler.typeEnforced)
	}
	// Collect inputs
	MainWorkLoop: for {
		select {
			case <-quitChan:
				close(assembler.assembledChan)
				return
			case posVal, ok := <-assembler.collectChan:
				if !ok {
					break MainWorkLoop
				}
				missing--
				results[posVal.position] = posVal.data
				if missing == 0 {
					break MainWorkLoop
				}
		}
	}
	assembler.assembledChan <- results
}

func assemblerWorker(quitChan executor.QuitChan, node Node, position int, collectChan collectChan, typeEnforced reflect.Type) {
	select {
		case <-quitChan:
			return
		// Read value to transmit
		case val, ok := <-node.Out():
			if !ok {
				return
			}
			if typeEnforced != nil && typeEnforced != reflect.TypeOf(val) {
				typeMismatchErrorForTypeOf(typeEnforced, val)
			}
			select {
				case <-quitChan:
					return
				// Transmit the value
				case collectChan <- PositionedData{position, val}:
			}
	}
}
