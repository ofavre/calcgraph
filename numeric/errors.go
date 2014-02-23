package numeric

import (
    "fmt"
    "reflect"
)

type NonNumericTypeError struct {
    runtimeType reflect.Type
}

func nonNumericTypeErrorForTypeOf(val interface{}) {
    nonNumericTypeError(reflect.TypeOf(val))
}

func nonNumericTypeError(runtimeType reflect.Type) {
    panic(NonNumericTypeError{runtimeType})
}

func (err NonNumericTypeError) Error() string {
    return fmt.Sprintf("%s is not a numeric type", err.runtimeType)
}
