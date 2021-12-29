package variable

type List []interface{}

func (l List) All() []Variable {
	var vars []Variable
	for _, v := range l {
		vars = append(vars, New(v))
	}
	return vars
}

func (l List) Len() int {
	return len(l)
}
