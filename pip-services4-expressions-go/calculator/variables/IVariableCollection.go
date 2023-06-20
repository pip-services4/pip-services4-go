package variables

// IVariableCollection defines a variables list.
type IVariableCollection interface {
	// Add a new variable to the collection.
	//	Parameters:
	//		- variable: a variable to be added.
	Add(variable IVariable)

	// Length number of variables stored in the collection.
	Length() int

	// Get a variable by its index.
	//	Parameters:
	//		- index: a variable index.
	//	Returns: a retrieved variable.
	Get(index int) IVariable

	// GetAll variables stores in the collection
	//	Returns: a list with variables.
	GetAll() []IVariable

	// FindIndexByName finds variable index in the list by it's name.
	//	Parameters:
	//		- name: The variable name to be found.
	//	Returns: Variable index in the list or <code>-1</code> if variable was not found.
	FindIndexByName(name string) int

	// FindByName finds variable in the list by it's name.
	//	Parameters:
	//		- name: The variable name to be found.
	//	Returns: Variable or <code>null</code> if function was not found.
	FindByName(name string) IVariable

	// Locate variable in the list or create a new one if variable was not found.
	//	Parameters:
	//		- name: The variable name to be found.
	//	Returns: Found or created variable.
	Locate(name string) IVariable

	// Remove a variable by its index.
	//	Parameters:
	//		- index: a index of the variable to be removed.
	Remove(index int)

	// RemoveByName removes variable by it's name.
	//	Parameters:
	//		- name: The variable name to be removed.
	RemoveByName(name string)

	// Clear the collection.
	Clear()

	// ClearValues clears all stored variables (assigns null values).
	ClearValues()
}
