package functions

import "strings"

// FunctionCollection implements a functions list.
type FunctionCollection struct {
	functions []IFunction
}

func NewFunctionCollection() *FunctionCollection {
	c := &FunctionCollection{
		functions: []IFunction{},
	}
	return c
}

// Add a new function to the collection.
//	Parameters:
//		- function: a function to be added.
func (c *FunctionCollection) Add(function IFunction) {
	if function == nil {
		panic("Function cannot be nil.")
	}
	c.functions = append(c.functions, function)
}

// Length is a number of functions stored in the collection.
func (c *FunctionCollection) Length() int {
	return len(c.functions)
}

// Get a function by its index.
//	Parameters:
//		- index: a function index.
//	Returns: a retrieved function.
func (c *FunctionCollection) Get(index int) IFunction {
	return c.functions[index]
}

// GetAll all functions stores in the collection
//	Returns: a list with functions.
func (c *FunctionCollection) GetAll() []IFunction {
	result := []IFunction{}
	result = append(result, c.functions...)
	return result
}

// FindIndexByName function index in the list by it's name.
//	Parameters:
//		- name: The function name to be found.
//	Returns: Function index in the list or <code>-1</code> if function was not found.
func (c *FunctionCollection) FindIndexByName(name string) int {
	name = strings.ToUpper(name)
	for i, f := range c.functions {
		if strings.ToUpper(f.Name()) == name {
			return i
		}
	}
	return -1
}

// FindByName finds function in the list by it's name.
//	Parameters:
//		- name: The function name to be found.
//	Returns: Function or <code>null</code> if function was not found.
func (c *FunctionCollection) FindByName(name string) IFunction {
	index := c.FindIndexByName(name)
	if index >= 0 {
		return c.functions[index]
	}
	return nil
}

// Remove a function by its index.
//	Parameters:
//		- index: a index of the function to be removed.
func (c *FunctionCollection) Remove(index int) {
	c.functions = append(c.functions[:index], c.functions[index+1:]...)
}

// RemoveByName function by it's name.
//	Parameters:
//		- name: The function name to be removed.
func (c *FunctionCollection) RemoveByName(name string) {
	index := c.FindIndexByName(name)
	if index >= 0 {
		c.functions = append(c.functions[:index], c.functions[index+1:]...)
	}
}

// Clear the collection.
func (c *FunctionCollection) Clear() {
	c.functions = []IFunction{}
}
