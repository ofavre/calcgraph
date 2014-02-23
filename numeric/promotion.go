package numeric

import (
    "reflect"
)

// Promote the arguments' type, so they are of the same type.
//
// Convert to the lowest type common that can contain the respective
// values of the arguments, without (or with least) precision loss.
//
// Note: that calling with uint8(0),int8(0) will return int16(0),int16(0),
// as the convertions are only based on the possible range of the
// runtime types of the arguments, not their runtime value.
func Promote(a, b interface{}) (interface{}, interface{}) {
    if reflect.TypeOf(a) == reflect.TypeOf(b) {
        return a, b
    }
    // Get rid of int and uint for their underlying native type
    a = Native(a)
    b = Native(b)
    // Convert to the lowest type that can contain the respective values
    // without (or with least) precision loss
    switch a := a.(type) {
        case int8:
            switch b := b.(type) {
                case int8:
                    return int8(a), int8(b)
                case uint8: //fallthrough
                    return int16(a), int16(b)
                case int16:
                    return int16(a), int16(b)
                case uint16: //fallthrough
                    return int32(a), int32(b)
                case int32:
                    return int32(a), int32(b)
                case uint32: //fallthrough
                    return int64(a), int64(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return float64(a), float64(b)
                case float32:
                    return float32(a), float32(b)
                case float64:
                    return float64(a), float64(b)
                case complex64:
                    return complex(float32(a), 0), complex64(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case uint8:
            switch b := b.(type) {
                case int8: // fallthrough
                    return int16(a), int16(b)
                case int16:
                    return int16(a), int16(b)
                case uint8:
                    return uint8(a), uint8(b)
                case uint16:
                    return uint16(a), uint16(b)
                case int32:
                    return int32(a), int32(b)
                case uint32:
                    return uint32(a), uint32(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return uint64(a), uint64(b)
                case float32:
                    return float32(a), float32(b)
                case float64:
                    return float64(a), float64(b)
                case complex64:
                    return complex(float32(a), 0), complex64(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case int16:
            switch b := b.(type) {
                case int8:  // fallthrough
                    return int16(a), int16(b)
                case uint8: // fallthrough
                    return int16(a), int16(b)
                case int16:
                    return int16(a), int16(b)
                case uint16: // fallthrough
                    return int32(a), int32(b)
                case int32:
                    return int32(a), int32(b)
                case uint32: // fallthrough
                    return int64(a), int64(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return float64(a), float64(b)
                case float32:
                    return float32(a), float32(b)
                case float64:
                    return float64(a), float64(b)
                case complex64:
                    return complex(float32(a), 0), complex64(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case uint16:
            switch b := b.(type) {
                case int8:  // fallthrough
                    return int32(a), int32(b)
                case int16: // fallthrough
                    return int32(a), int32(b)
                case int32:
                    return int32(a), int32(b)
                case uint8:  // fallthrough
                    return uint32(a), uint32(b)
                case uint16: // fallthrough
                    return uint32(a), uint32(b)
                case uint32:
                    return uint32(a), uint32(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return uint64(a), uint64(b)
                case float32:
                    return float32(a), float32(b)
                case float64:
                    return float64(a), float64(b)
                case complex64:
                    return complex(float32(a), 0), complex64(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case int32:
            switch b := b.(type) {
                case int8:   // fallthrough
                    return int32(a), int32(b)
                case uint8:  // fallthrough
                    return int32(a), int32(b)
                case int16:  // fallthrough
                    return int32(a), int32(b)
                case uint16: // fallthrough
                    return int32(a), int32(b)
                case int32:
                    return int32(a), int32(b)
                case uint32:
                    return int64(a), int64(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return float64(a), float64(b)
                case float32: // fallthrough
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64: // fallthrough
                    return complex(float64(a), 0), complex128(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case uint32:
            switch b := b.(type) {
                case uint8:  // fallthrough
                    return uint32(a), uint32(b)
                case uint16: // fallthrough
                    return uint32(a), uint32(b)
                case uint32:
                    return uint32(a), uint32(b)
                case uint64:
                    return uint64(a), uint64(b)
                case int8:  // fallthrough
                    return int64(a), int64(b)
                case int16: // fallthrough
                    return int64(a), int64(b)
                case int32: // fallthrough
                    return int64(a), int64(b)
                case int64:
                    return int64(a), int64(b)
                case float32: // fallthrough
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64: // fallthrough
                    return complex(float64(a), 0), complex128(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case int64:
            switch b := b.(type) {
                case int8:   // fallthrough
                    return int64(a), int64(b)
                case uint8:  // fallthrough
                    return int64(a), int64(b)
                case int16:  // fallthrough
                    return int64(a), int64(b)
                case uint16: // fallthrough
                    return int64(a), int64(b)
                case int32:  // fallthrough
                    return int64(a), int64(b)
                case uint32: // fallthrough
                    return int64(a), int64(b)
                case int64:
                    return int64(a), int64(b)
                case uint64:
                    return float64(a), float64(b)
                case float32: // fallthrough
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64: // fallthrough
                    return complex(float64(a), 0), complex128(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case uint64:
            switch b := b.(type) {
                case uint8:  // fallthrough
                    return uint64(a), uint64(b)
                case uint16: // fallthrough
                    return uint64(a), uint64(b)
                case uint32: // fallthrough
                    return uint64(a), uint64(b)
                case uint64:
                    return uint64(a), uint64(b)
                case int8:    // fallthrough
                    return float64(a), float64(b)
                case int16:   // fallthrough
                    return float64(a), float64(b)
                case int32:   // fallthrough
                    return float64(a), float64(b)
                case int64:   // fallthrough
                    return float64(a), float64(b)
                case float32:
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64: // fallthrough
                    return complex(float64(a), 0), complex128(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case float32:
            switch b := b.(type) {
                case int8:   // fallthrough
                    return float32(a), float32(b)
                case uint8:  // fallthrough
                    return float32(a), float32(b)
                case int16:  // fallthrough
                    return float32(a), float32(b)
                case uint16:
                    return float32(a), float32(b)
                case int32:   // fallthrough
                    return float64(a), float64(b)
                case uint32:  // fallthrough
                    return float64(a), float64(b)
                case int64:   // fallthrough
                    return float64(a), float64(b)
                case uint64:  // fallthrough
                    return float64(a), float64(b)
                case float32: // fallthrough
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64:
                    return complex(float32(a), 0), complex64(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case float64:
            switch b := b.(type) {
                case int8:    // fallthrough
                    return float64(a), float64(b)
                case uint8:   // fallthrough
                    return float64(a), float64(b)
                case int16:   // fallthrough
                    return float64(a), float64(b)
                case uint16:  // fallthrough
                    return float64(a), float64(b)
                case int32:   // fallthrough
                    return float64(a), float64(b)
                case uint32:  // fallthrough
                    return float64(a), float64(b)
                case int64:   // fallthrough
                    return float64(a), float64(b)
                case uint64:  // fallthrough
                    return float64(a), float64(b)
                case float32: // fallthrough
                    return float64(a), float64(b)
                case float64:
                    return float64(a), float64(b)
                case complex64: // fallthrough
                    return complex(float64(a), 0), complex128(b)
                case complex128:
                    return complex(float64(a), 0), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case complex64:
            switch b := b.(type) {
                case int8:   // fallthrough
                    return complex64(a), complex(float32(b), 0)
                case uint8:  // fallthrough
                    return complex64(a), complex(float32(b), 0)
                case int16:  // fallthrough
                    return complex64(a), complex(float32(b), 0)
                case uint16:
                    return complex64(a), complex(float32(b), 0)
                case int32:   // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint32:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case int64:   // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint64:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case float32: // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case float64:
                    return complex128(a), complex(float64(b), 0)
                case complex64:
                    return complex64(a), complex64(b)
                case complex128:
                    return complex128(a), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        case complex128:
            switch b := b.(type) {
                case int8:   // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint8:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case int16:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint16:
                    return complex128(a), complex(float64(b), 0)
                case int32:   // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint32:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case int64:   // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case uint64:  // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case float32: // fallthrough
                    return complex128(a), complex(float64(b), 0)
                case float64:
                    return complex128(a), complex(float64(b), 0)
                case complex64: // fallthrough
                    return complex128(a), complex128(b)
                case complex128:
                    return complex128(a), complex128(b)
                default:
                    nonNumericTypeErrorForTypeOf(b)
            }
        default:
            nonNumericTypeErrorForTypeOf(a)
    }
    panic("reached unreachable point")
}
