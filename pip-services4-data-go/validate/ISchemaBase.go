package validate

type ISchemaBase interface {
	PerformValidation(path string, value any) []*ValidationResult
}
