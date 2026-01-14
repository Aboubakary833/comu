package validation

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/utils"
	"regexp"

	"github.com/Oudwins/zog"
)

var (
	msgNameRequired 		 		= "Name is required"
	mgsEmailRequired 		 		= "Email address is required"
	mgsPasswordRequired 	 		= "Password is required"
	msgInvalidEmail			 		= "Provided email is invalid"
	msgNameTooBig			 		= "Name must not be more than 50 characters long"
	msgNameTooShort			 		= "Name must be at least 3 characters long"
	msgPasswordMustHaveDigit 		= "Password must contain at least one digit"
	msgPasswordMustHaveUpperCase 	= "Password must contain at least one uppercase letter"
	msgPasswordMustHaveSpecialChar  = "Password must contain at least one special character"
	mgsInvalidOtp 			 		= utils.UcFirst(domain.ErrInvalidOtp.Error())
)


type StructValidator struct {
	shape *zog.StructSchema
}

type SchemaValidationErrors map[string]string

func NewStructValidator(shape *zog.StructSchema) *StructValidator {
	return &StructValidator{
		shape: shape,
	}
}

func (validator *StructValidator) Validate(data, destPtr any, opt ...zog.ExecOption) SchemaValidationErrors {
	errList := validator.shape.Parse(data, destPtr, opt...)
	
	if errList != nil {
		var errors = make(SchemaValidationErrors)

		for _, err := range errList {
			errors[err.PathString()] = err.Message
		}

		return errors
	}

	return nil
}

var LoginValidator = NewStructValidator(zog.Struct(zog.Shape{
	"email": zog.String().Required(zog.Message(mgsEmailRequired)),
	"password": zog.String().Required(zog.Message(mgsPasswordRequired)),
}))

var RegisterValidator = NewStructValidator(zog.Struct(zog.Shape{
	"name": zog.String().Required(zog.Message(msgNameRequired)).
	Min(3, zog.Message(msgNameTooShort)).Max(50, zog.Message(msgNameTooBig)),
	"email": zog.String().Required(zog.Message(mgsEmailRequired)).Email(zog.Message(msgInvalidEmail)),
	"password": zog.String().Required(zog.Message(mgsPasswordRequired)).
	ContainsDigit(zog.Message(msgPasswordMustHaveDigit)).
	ContainsUpper(zog.Message(msgPasswordMustHaveUpperCase)).
	ContainsSpecial(zog.Message(msgPasswordMustHaveSpecialChar)),
}))

var OtpCodeValidator = NewStructValidator(zog.Struct(zog.Shape{
	"email": zog.String().Required(zog.Message(mgsEmailRequired)),
	"code": zog.String().Len(6, zog.Message(mgsInvalidOtp)).
	Match(regexp.MustCompile("[0-9]"), zog.Message(mgsInvalidOtp)),
}))

