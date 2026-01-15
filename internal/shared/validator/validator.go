package validator

import "github.com/Oudwins/zog"


type StructValidator struct {
	schema *zog.StructSchema
}

type SchemaValidationErrors map[string]string

func NewStructValidator(zogStruct *zog.StructSchema) *StructValidator {
	return &StructValidator{
		schema: zogStruct,
	}
}

func (validator *StructValidator) Validate(data, destPtr any, opt ...zog.ExecOption) SchemaValidationErrors {
	errList := validator.schema.Parse(data, destPtr, opt...)
	
	if errList != nil {
		var errors = make(SchemaValidationErrors)

		for _, err := range errList {
			errors[err.PathString()] = err.Message
		}

		return errors
	}

	return nil
}
