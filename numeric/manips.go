package numeric

import (
    "fmt"
    "unsafe"
    "reflect"
)

func init() {
    var i int
    switch l := unsafe.Sizeof(i) ; l {
        case 4, 8:
        default:
            panic(fmt.Sprintf("unknown Sizeof(int) = %d", l))
    }
}

func ZeroOfSameType(val interface{}) interface{} {
	return ZeroOfType(reflect.TypeOf(val))
}

func ZeroOfType(typ reflect.Type) interface{} {
	return reflect.Zero(typ).Interface()
}

func Native(val interface{}) interface{} {
    switch val := val.(type) {
        case int:
            switch unsafe.Sizeof(val) {
                case 4:
                    return int32(val)
                case 8:
                    return int64(val)
            }
        case uint:
            switch unsafe.Sizeof(val) {
                case 4:
                    return uint32(val)
                case 8:
                    return uint64(val)
            }
        default:
    }
    return val
}
