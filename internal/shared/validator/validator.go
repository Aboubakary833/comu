package validator

import (
	"github.com/Oudwins/zog"
)

type StructValidator struct {
	schema *zog.StructSchema
}

type SchemaValidationErrors map[string]string

func NewStructValidator(zogStruct *zog.StructSchema) *StructValidator {
	return &StructValidator{
		schema: zogStruct,
	}
}

func (validator *StructValidator) Validate(dataPtr any, opt ...zog.ExecOption) SchemaValidationErrors {
	errList := validator.schema.Validate(dataPtr, opt...)

	if errList != nil {
		var errors = make(SchemaValidationErrors)

		for _, err := range errList {
			path := err.PathString()
			//TODO: check if path is camelcase and transform it to snake_case
			errors[path] = err.Message
		}

		return errors
	}

	return nil
}
