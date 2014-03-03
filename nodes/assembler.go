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



type AssemblerNode struct {
	out				DataChan
	typeEnforced	reflect.Type
	inNodes			[]Node
	selectCases		[]reflect.SelectCase
	remaining		int
	lastQuitChan	executor.QuitChan
}

func NewAssemblerNode(typeEnforced reflect.Type, inNodes ...Node) *AssemblerNode {
	// Construct the select cases from each node
	selectCases := make([]reflect.SelectCase, 1+len(inNodes)) // first item is reserved for the QuitChan
	var nilValue reflect.Value
	for i, node := range inNodes {
		selectCases[1+i] = reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(node.Out()), nilValue}
	}
	return &AssemblerNode{make(DataChan), typeEnforced, inNodes, selectCases, len(inNodes), nil}
}

func (node AssemblerNode) Out() DataChan {
	return node.out
}

func (node *AssemblerNode) Run(quitChan executor.QuitChan) {
	missing	:= node.remaining
	results	:= make([]Data, len(node.inNodes))
	var nilChanValue reflect.Value = reflect.ValueOf(nil)
	// Update the QuitChan select case if needed
	if (node.lastQuitChan != quitChan) {
		var nilValue reflect.Value
		node.selectCases[0] = reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(quitChan), nilValue}
	}
	// Copy every select cases, they'll be nil-ed out one by one
	selectCases := make([]reflect.SelectCase, len(node.selectCases))
	copy(selectCases, node.selectCases[:])
	// Read each input chan one by one
	for missing > 0 && node.remaining > 0 {
		chosen, val, ok := reflect.Select(selectCases)
		if chosen == 0 { // quitChan
			close(node.out)
			return
		} else {
			// Got a value from one chan
			missing--
			// Don't read this chan again in this call to Run()
			selectCases[chosen].Chan = nilChanValue
			// Type enforcement, only if we got a real value, and not a nil for chan closed
			if ok && node.typeEnforced != nil && node.typeEnforced != reflect.TypeOf(val) {
				typeMismatchErrorForTypeOf(node.typeEnforced, val)
			}
			// Set the corresponding output value
			results[chosen-1] = Data(val.Interface())
			// If we got a nil for chan closed, remove the chan
			if !ok {
				// We remove it from the object's selectCases,
				// not to reread it in a future call to Run()
				node.remaining--
				node.selectCases[chosen].Chan = nilChanValue
			}
		}
	}
	// Return the assembled values
	if missing <= 0 {
		select {
			case <-quitChan:
				close(node.out)
				return
			case node.out <- results:
		}
	}
	// Close if all our input nodes are closed
	if node.remaining <= 0 {
		close(node.out)
		<-quitChan
		return
	}
}
