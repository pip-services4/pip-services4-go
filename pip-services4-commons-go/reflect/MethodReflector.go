package reflect

import (
	refl "reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MethodReflector helper class to perform method introspection and dynamic invocation.
//
// This class has symmetric implementation across all languages supported by Pip.Services toolkit
// and used to support dynamic data processing.
//
// Because all languages have different casing and case sensitivity rules,
// this MethodReflector treats all method names as case insensitive.
//
//	Example:
//		myObj = MyObject();
//
//		methods = MethodReflector.GetMethodNames();
//		MethodReflector.HasMethod(myObj, "myMethod");
//		MethodReflector.InvokeMethod(myObj, "myMethod", 123);
var MethodReflector = &_TMethodReflector{}

type _TMethodReflector struct{}

func (c *_TMethodReflector) matchMethod(method refl.Method, name string) bool {
	// Method must be public and match to name as case insensitive
	r, _ := utf8.DecodeRuneInString(method.Name)
	return unicode.IsUpper(r) &&
		strings.ToLower(method.Name) == strings.ToLower(name)
}

// HasMethod checks if object has a method with specified name..
//	Parameters:
//		- obj any an object to introspect.
//		- name string a name of the method to check.
//	Returns: bool true if the object has the method and false if it doesn't.
func (c *_TMethodReflector) HasMethod(obj any, name string) bool {
	if obj == nil {
		panic("Object cannot be nil")
	}
	if name == "" {
		panic("Method name cannot be empty")
	}

	objType := refl.TypeOf(obj)

	for index := 0; index < objType.NumMethod(); index++ {
		method := objType.Method(index)
		if c.matchMethod(method, name) {
			return true
		}
	}

	return false
}

// InvokeMethod invokes an object method by its name with specified parameters.
//	Parameters:
//		- obj any an object to invoke.
//		- name string a name of the method to invoke.
//		- args ...any a list of method arguments.
//	Returns: any the result of the method invocation or null if method returns void.
func (c *_TMethodReflector) InvokeMethod(obj any, name string, args ...any) any {
	if obj == nil {
		panic("Object cannot be nil")
	}
	if name == "" {
		panic("Method name cannot be empty")
	}

	defer func() {
		// Do nothing and return nil
		recover()
	}()

	objType := refl.TypeOf(obj)

	// Convert arguments
	inputs := make([]refl.Value, len(args))
	for index := range args {
		inputs[index] = refl.ValueOf(args[index])
	}

	for index := 0; index < objType.NumMethod(); index++ {
		method := objType.Method(index)
		if c.matchMethod(method, name) {
			results := refl.ValueOf(obj).Method(index).Call(inputs)
			if len(results) == 0 {
				return nil
			}
			return results[0].Interface()
		}
	}

	return nil
}

// GetMethodNames gets names of all methods implemented in specified object.
// Parameters:
//  - obj any
//  an objec to introspect.
// Returns []string
// a list with method names.
func (c *_TMethodReflector) GetMethodNames(obj any) []string {
	objType := refl.TypeOf(obj)
	methods := make([]string, 0, objType.NumMethod())

	for index := 0; index < objType.NumMethod(); index++ {
		method := objType.Method(index)
		if c.matchMethod(method, method.Name) {
			methods = append(methods, method.Name)
		}
	}

	return methods
}
