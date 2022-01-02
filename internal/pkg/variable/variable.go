package variable

import (
	"fmt"
	"strconv"
	"strings"
)

type Variable interface {
	AsString() string
	AsInteger() int
	AsFloat64() float64
	AsBool() bool
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

func (v *variable) AsBool() bool {
	switch actual := v.raw.(type) {
	case bool:
		return actual
	case string:
		return actual != ""
	case int:
		return actual > 0
	case float64:
		return actual > 0
	}
	return false
}

func (v *variable) AsString() string {
	switch actual := v.raw.(type) {
	case bool:
		return fmt.Sprintf("%t", actual)
	case string:
		return actual
	case int:
		return fmt.Sprintf("%d", actual)
	case float64:
		return fmt.Sprintf("%f", actual)
	case []int:
		var out []string
		for _, a := range actual {
			out = append(out, fmt.Sprintf("%d", a))
		}
		return fmt.Sprintf("[%s]", strings.Join(out, ", "))
	case []float64:
		var out []string
		for _, a := range actual {
			out = append(out, fmt.Sprintf("%f", a))
		}
		return fmt.Sprintf("[%s]", strings.Join(out, ", "))
	case []string:
		return fmt.Sprintf(`["%s"]`, strings.Join(actual, `", "`))
	}
	return fmt.Sprintf("%v", v.raw)
}

func (v *variable) AsInteger() int {
	switch c := v.raw.(type) {
	case int:
		return c
	case float64:
		return int(c)
	case bool:
		if c {
			return 1
		}
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
	case bool:
		if c {
			return 1
		}
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
	switch actual := v.raw.(type) {
	case []interface{}:
		return List(actual)
	case []string:
		var anon []interface{}
		for _, s := range actual {
			anon = append(anon, s)
		}
		return List(anon)
	case []float64:
		var anon []interface{}
		for _, s := range actual {
			anon = append(anon, s)
		}
		return List(anon)
	case []bool:
		var anon []interface{}
		for _, s := range actual {
			anon = append(anon, s)
		}
		return List(anon)
	case []int:
		var anon []interface{}
		for _, s := range actual {
			anon = append(anon, s)
		}
		return List(anon)
	}
	return List([]interface{}{})
}
