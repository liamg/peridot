package variable

import (
	"fmt"
	"strconv"
)

type Variable interface {
	AsString() string
	AsInteger() int
	AsFloat64() float64
	AsCollection() Collection
	AsList() List
	Interface() interface{}
}

type variable struct {
	raw interface{}
}

func New(raw interface{}) Variable {
	return &variable{
		raw: raw,
	}
}

func (v *variable) Interface() interface{} {
	return v.raw
}

func (v *variable) AsString() string {
	return fmt.Sprintf("%s", v.raw)
}

func (v *variable) AsInteger() int {
	switch c := v.raw.(type) {
	case int:
		return c
	case float64:
		return int(c)
	case string:
		if i, err := strconv.Atoi(c); err == nil {
			return i
		}
	}
	return 0
}

func (v *variable) AsFloat64() float64 {
	switch c := v.raw.(type) {
	case int:
		return float64(c)
	case float64:
		return c
	case string:
		if f, err := strconv.ParseFloat(c, 64); err == nil {
			return f
		}
	}
	return 0
}

func (v *variable) AsCollection() Collection {
	if m, ok := v.raw.(map[string]interface{}); ok {
		return NewCollection(m)
	}
	return NewCollection(nil)
}

func (v *variable) AsList() List {
	if m, ok := v.raw.([]interface{}); ok {
		return List(m)
	}
	return List([]interface{}{})
}
