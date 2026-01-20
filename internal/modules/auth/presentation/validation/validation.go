package validation

import (
	"comu/internal/modules/auth/domain"
	"comu/internal/shared/utils"
	"comu/internal/shared/validator"
	"regexp"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/internals"
)

var (
	msgNameRequired                = "Name is required"
	msgEmailRequired               = "Email address is required"
	msgPasswordRequired            = "Password is required"
	msgPasswordShouldBeConfirmed   = "Password should be confirmed"
	msgPasswordsDoNotMatch         = "Password and password confirmation don't match"
	msgTokenRequired               = "Token is required"
	msgInvalidEmail                = "Provided email is invalid"
	msgNameTooBig                  = "Name must not be more than 50 characters long"
	msgNameTooShort                = "Name must be at least 3 characters long"
	msgPasswordMustHaveDigit       = "Password must contain at least one digit"
	msgPasswordMustHaveUpperCase   = "Password must contain at least one uppercase letter"
	msgPasswordMustHaveSpecialChar = "Password must contain at least one special character"
	msgInvalidOtp                  = utils.UcFirst(domain.ErrInvalidOtp.Error())
)

var confirmed = zog.CustomFunc(func(ptr *string, ctx zog.Ctx) bool {
	password, ok := ctx.Get("password").(string)
	if !ok {
		return false
	}

	return *ptr == password
})

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
	"reset_token": zog.String().Required(zog.Message(msgTokenRequired)),
	"password":    zog.String().Required(zog.Message(msgPasswordRequired)),
	"passwordConfirmation": zog.String().Required(zog.Message(msgPasswordShouldBeConfirmed)).TestFunc(func(val *string, ctx internals.Ctx) bool {
		password, ok := ctx.Get("password").(string)
		if !ok || *val == password {
			ctx.AddIssue(&zog.ZogIssue{
				Message: "Password and password confirmation don't match",
				Path:    []string{"passwordConfirmation"},
				Value:   val,
			})
		}

		return false
	}),
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
	"email":        zog.String().Required(zog.Message(msgEmailRequired)).Email(zog.Message(msgInvalidEmail)),
	"resendToken": zog.String().Required(zog.Message(msgTokenRequired)),
}))
