package validation

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/utils"
	"comu/internal/shared/validator"
	"regexp"

	"github.com/Oudwins/zog"
)

var (
	msgNameRequired                = "Name is required"
	msgEmailRequired               = "Email address is required"
	msgPasswordRequired            = "Password is required"
	msgPasswordShouldBeConfirmed   = "Password should be confirmed"
	msgTokenRequired               = "Token is required"
	msgInvalidEmail                = "Provided email is invalid"
	msgNameTooBig                  = "Name must not be more than 50 characters long"
	msgNameTooShort                = "Name must be at least 3 characters long"
	msgPasswordMustHaveDigit       = "Password must contain at least one digit"
	msgPasswordMustHaveUpperCase   = "Password must contain at least one uppercase letter"
	msgPasswordMustHaveSpecialChar = "Password must contain at least one special character"
	msgInvalidOtp                  = utils.UcFirst(domain.ErrInvalidOtp.Error())
)

var LoginValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"email":    zog.String().Required(zog.Message(msgEmailRequired)),
	"password": zog.String().Required(zog.Message(msgPasswordRequired)),
}))

var RegisterValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"name": zog.String().Required(zog.Message(msgNameRequired)).
		Min(3, zog.Message(msgNameTooShort)).Max(50, zog.Message(msgNameTooBig)),
	"email": zog.String().Required(zog.Message(msgEmailRequired)).Email(zog.Message(msgInvalidEmail)),
	"password": zog.String().Required(zog.Message(msgPasswordRequired)).
		ContainsDigit(zog.Message(msgPasswordMustHaveDigit)).
		ContainsUpper(zog.Message(msgPasswordMustHaveUpperCase)).
		ContainsSpecial(zog.Message(msgPasswordMustHaveSpecialChar)),
}))

var NewPasswordValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"resetToken":           zog.String().Required(zog.Message(msgTokenRequired)),
	"password":             zog.String().Required(zog.Message(msgPasswordRequired)),
	"passwordConfirmation": zog.String().Required(zog.Message(msgPasswordShouldBeConfirmed)),
}))

var ResetPasswordValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"email": zog.String().Required(zog.Message(msgEmailRequired)).Email(zog.Message(msgInvalidEmail)),
}))

var OtpCodeValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"email": zog.String().Required(zog.Message(msgEmailRequired)),
	"code": zog.String().Len(6, zog.Message(msgInvalidOtp)).
		Match(regexp.MustCompile("[0-9]"), zog.Message(msgInvalidOtp)),
}))

var ResendOtpValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"email":       zog.String().Required(zog.Message(msgEmailRequired)).Email(zog.Message(msgInvalidEmail)),
	"resendToken": zog.String().Required(zog.Message(msgTokenRequired)),
}))
