package variable

type Collection interface {
	Has(name string) bool
	Get(name string) Variable
	Set(name string, value interface{})
	AsMap() map[string]interface{}
	MergeIn(Collection)
}

func NewCollection(data map[string]interface{}) Collection {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &collection{
		data: data,
	}
}

type collection struct {
	data map[string]interface{}
}

func (c *collection) Has(name string) bool {
	_, ok := c.data[name]
	return ok
}

func (c *collection) Get(name string) Variable {
	if value, ok := c.data[name]; ok {
		return New(value)
	}
	return New(nil)
}

func (c *collection) Set(name string, value interface{}) {
	c.data[name] = value
}

func (c *collection) AsMap() map[string]interface{} {
	return c.data
}

// MergeIn takes collection c2 and writes all of it's values to c, regardless of whether they already exist
func (c *collection) MergeIn(c2 Collection) {
	for k, v := range c2.AsMap() {
		c.data[k] = v
	}
}
