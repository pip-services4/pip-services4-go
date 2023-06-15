package data

// MultiString An object that contains string translations for multiple languages.
// Language keys use two-letter codes like: 'en', 'sp', 'de', 'ru', 'fr', 'pr'.
// When translation for specified language does not exists it defaults to English ('en').
// When English does not exists it falls back to the first defined language.
//
//	Example:
//	values := NewMultiStringFromTuples(
//		"en", "Hello World!",
//		"ru", "Привет мир!"
//	);
//
//	value1 := values.Get("ru"); // Result: "Привет мир!"
//	value2 := values.Get("pt"); // Result: "Hello World!"
type MultiString map[string]string
