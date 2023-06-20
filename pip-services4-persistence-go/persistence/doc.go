// Package persistence contains interfaces for various design patterns that work with data.
// Contains various persistence implementations (InMemory and File –persistences).
// These are "abstract" persistences, which only connect to data sources and do not implement
// the operations and methods for working the data. The classes that extend these persistences
// must implement this logic on their own.
// Identifiable Persistences work with Identifiable objects, which have primary keys.
// A few standard operations are defined by default for these objects:
// reading arrays and data pages; searching for an object by its id;
// and creating, updating, and deleting records of objects.

package persistence
