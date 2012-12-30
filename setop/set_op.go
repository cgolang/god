package setop

import (
	"bytes"
	"fmt"
	"strings"
)

type RawSourceCreator func(b []byte) Skipper
type SetOpResultIterator func(res *SetOpResult)

const (
	Append = iota
	ConCat
	IntegerSum
	IntegerSub
	IntegerDiv
	IntegerMul
	FloatSum
	FloatSub
	FloatDiv
	FloatMul
	BigIntAnd
	BigIntAdd
	BigIntAndNot
	BigIntDiv
	BigIntMod
	BigIntMul
	BigIntOr
	BigIntRem
	BigIntSub
	BigIntXor
)

type SetOpMerge int

func ParseSetOpMerge(s string) (result SetOpMerge, err error) {
	switch s {
	case "Append":
		result = Append
	case "ConCat":
		result = ConCat
	case "IntegerSum":
		result = IntegerSum
	case "IntegerSub":
		result = IntegerSub
	case "IntegerDiv":
		result = IntegerDiv
	case "IntegerMul":
		result = IntegerMul
	case "FloatSum":
		result = FloatSum
	case "FloatSub":
		result = FloatSub
	case "FloatDiv":
		result = FloatDiv
	case "FloatMul":
		result = FloatMul
	case "BigIntAnd":
		result = BigIntAnd
	case "BigIntAdd":
		result = BigIntAdd
	case "BigIntAndNot":
		result = BigIntAndNot
	case "BigIntDiv":
		result = BigIntDiv
	case "BigIntMod":
		result = BigIntMod
	case "BigIntMul":
		result = BigIntMul
	case "BigIntOr":
		result = BigIntOr
	case "BigIntRem":
		result = BigIntRem
	case "BigIntSub":
		result = BigIntSub
	case "BigIntXor":
		result = BigIntXor
	default:
		err = fmt.Errorf("Unknown SetOpType %v. Legal values: Append, ConCat, IntegerSum, IntegerSub, IntegerDiv, IntegerMul, FloatSum, FloatSub, FloatDiv, FloatMul, BigIntAdd, BigIntAnd, BigIntAndNot, BigIntDiv, BigIntMod, BigIntMul, BigIntOr, BigIntRem, BigIntSub, BigIntXor.", s)
	}
	return
}

func (self SetOpMerge) String() string {
	switch self {
	case Append:
		return "Append"
	case ConCat:
		return "ConCat"
	case IntegerSum:
		return "IntegerSum"
	case IntegerSub:
		return "IntegerSub"
	case IntegerDiv:
		return "IntegerDiv"
	case IntegerMul:
		return "IntegerMul"
	case FloatSum:
		return "FloatSum"
	case FloatSub:
		return "FloatSub"
	case FloatDiv:
		return "FloatDiv"
	case FloatMul:
		return "FloatMul"
	case BigIntAnd:
		return "BigIntAnd"
	case BigIntAdd:
		return "BigIntAdd"
	case BigIntAndNot:
		return "BigIntAndNot"
	case BigIntDiv:
		return "BigIntDiv"
	case BigIntMod:
		return "BigIntMod"
	case BigIntMul:
		return "BigIntMul"
	case BigIntOr:
		return "BigIntOr"
	case BigIntRem:
		return "BigIntRem"
	case BigIntSub:
		return "BigIntSub"
	case BigIntXor:
		return "BigIntXor"
	}
	panic(fmt.Errorf("Unknown SetOpType %v", int(self)))
}

const (
	Union = iota
	Intersection
	Difference
	Xor
)

type SetOpType int

func (self SetOpType) String() string {
	switch self {
	case Union:
		return "U"
	case Intersection:
		return "I"
	case Difference:
		return "D"
	case Xor:
		return "X"
	}
	panic(fmt.Errorf("Unknown SetOpType %v", int(self)))
}

type SetOpSource struct {
	Key   []byte
	SetOp *SetOp
}

type SetOp struct {
	Sources []SetOpSource
	Type    SetOpType
	Merge   SetOpMerge
}

func (self *SetOp) String() string {
	sources := make([]string, len(self.Sources))
	for index, source := range self.Sources {
		if source.Key != nil {
			sources[index] = string(source.Key)
		} else {
			sources[index] = fmt.Sprint(source.SetOp)
		}
	}
	return fmt.Sprintf("(%v %v)", self.Type, strings.Join(sources, " "))
}

type SetExpression struct {
	Op     *SetOp
	Code   string
	Min    []byte
	Max    []byte
	MinInc bool
	MaxInc bool
	Len    int
	Dest   []byte
}

func (self *SetExpression) Each(r RawSourceCreator, f SetOpResultIterator) (err error) {
	if self.Op == nil {
		self.Op = MustParse(self.Code)
	}
	skipper := createSkipper(r, self.Op)
	min := self.Min
	mininc := self.MinInc
	count := 0
	gt := -1
	if self.MaxInc {
		gt = 0
	}
	var res *SetOpResult
	for res, err = skipper.Skip(min, mininc); res != nil && err == nil; res, err = skipper.Skip(min, mininc) {
		if (self.Len > 0 && count >= self.Len) || (self.Max != nil && bytes.Compare(res.Key, self.Max) > gt) {
			return
		}
		count++
		min = res.Key
		mininc = false
		f(res)
	}
	return
}

type SetOpResult struct {
	Key    []byte
	Values [][]byte
}

func (self *SetOpResult) ShallowCopy() (result *SetOpResult) {
	result = &SetOpResult{
		Key:    self.Key,
		Values: make([][]byte, len(self.Values)),
	}
	copy(result.Values, self.Values)
	return
}
func (self *SetOpResult) String() string {
	return fmt.Sprintf("%+v", *self)
}