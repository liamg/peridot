package validate

import "fmt"

type Validator struct {
	vars map[string]interface{}
}

func New(vars map[string]interface{}) *Validator {
	return &Validator{
		vars: vars,
	}
}

func (v *Validator) EnsureDefined(name string) error {
	if _, ok := v.vars[name]; ok {
		return nil
	}
	return fmt.Errorf("required variable '%s' is not defined", name)
}

func (v *Validator) EnsureStringIfDefined(name string) error {
	if _, ok := v.vars[name]; !ok {
		return nil
	}
	return v.EnsureString(name)
}

func (v *Validator) EnsureStrings(names ...string) error {
	for _, name := range names {
		if err := v.EnsureString(name); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) EnsureString(name string) error {
	if err := v.EnsureDefined(name); err != nil {
		return err
	}
	value := v.vars[name]
	if _, ok := value.(string); ok {
		return nil
	}
	return fmt.Errorf("variable '%s' must be a string", name)
}

func (v *Validator) EnsureStringSliceIfDefined(name string) error {
	if _, ok := v.vars[name]; !ok {
		return nil
	}
	return v.EnsureStringSlice(name)
}

func (v *Validator) EnsureStringSlice(name string) error {
	if err := v.EnsureDefined(name); err != nil {
		return err
	}
	value := v.vars[name]
	if list, ok := value.([]interface{}); ok {
		if len(list) == 0 {
			return nil
		}
		for _, item := range list {
			if _, ok := item.(string); ok {
				return nil
			}
		}
	}
	return fmt.Errorf("variable '%s' must be a list of strings", name)
}
