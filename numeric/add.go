package numeric

func Add(values ...interface{}) interface{} {
	sum := values[0]
	for _, val := range(values[1:]) {
		sum = add(sum, val)
	}
	return sum
}

func add(_a, _b interface{}) interface{} {
	a, b := Promote(_a, _b)
	switch a := a.(type) {
		case int8:
			return a + b.(int8)
		case uint8:
			return a + b.(uint8)
		case int16:
			return a + b.(int16)
		case uint16:
			return a + b.(uint16)
		case int:
			return a + b.(int)
		case uint:
			return a + b.(uint)
		case int32:
			return a + b.(int32)
		case uint32:
			return a + b.(uint32)
		case int64:
			return a + b.(int64)
		case uint64:
			return a + b.(uint64)
		case float32:
			return a + b.(float32)
		case float64:
			return a + b.(float64)
		case complex64:
			return a + b.(complex64)
		case complex128:
			return a + b.(complex128)
		default:
			nonNumericTypeErrorForTypeOf(a)
	}
    panic("reached unreachable point")
}
