package css

import "fmt"

// CSS numeric value. Stored as either int64 or float64
type Num struct {
	Type  NumType // Type of number
	Value any     // Actual value of the number (either int64 or float64)
}

type NumType uint8

const (
	NumTypeInt = NumType(iota)
	NumTypeFloat
)

func (n Num) ToInt() int64 {
	if n.Type == NumTypeFloat {
		return int64(n.ToFloat())
	} else {
		return n.Value.(int64)
	}
}
func (n Num) ToFloat() float64 {
	if n.Type == NumTypeInt {
		return float64(n.ToInt())
	} else {
		return n.Value.(float64)
	}
}

func (n Num) Equals(other Num) bool {
	isFloat := (n.Type == NumTypeFloat) || (other.Type == NumTypeFloat)
	if isFloat {
		return n.ToFloat() == other.ToFloat()
	} else {
		return n.ToInt() == other.ToInt()
	}
}

func (n Num) Clamp(min, max Num) Num {
	// Using float all the time is probably fine, but let's avoid it if we can.
	isFloat := (n.Type == NumTypeFloat) || (min.Type == NumTypeFloat) || (max.Type == NumTypeFloat)
	if isFloat {
		if n.ToFloat() < min.ToFloat() {
			return min
		} else if max.ToFloat() < n.ToFloat() {
			return max
		}
	} else {
		if n.ToInt() < min.ToInt() {
			return min
		} else if max.ToInt() < n.ToInt() {
			return max
		}
	}
	return n
}

func (n Num) String() string {
	if n.Type == NumTypeFloat {
		return fmt.Sprintf("%f", n.ToFloat())
	}
	return fmt.Sprintf("%d", n.ToInt())
}

func NumFromInt(v int64) Num {
	return Num{NumTypeInt, v}
}
func NumFromFloat(v float64) Num {
	return Num{NumTypeFloat, v}
}
