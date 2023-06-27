package controllers

// ISwaggerController Interface to perform Swagger registrations.
type ISwaggerController interface {

	// RegisterOpenApiSpec Perform required Swagger registration steps.
	RegisterOpenApiSpec(baseRoute string, swaggerRoute string)
}
