package validation

import (
	"comu/internal/shared/validator"

	"github.com/Oudwins/zog"
)

var (
	msgPostIdRequired         = "Post id is required"
	msgTitleRequired          = "Title is required"
	msgContentRequired        = "Content required"
	msgTitleTooShort          = "Title must be more than 4 charaters long"
	msgTitleTooLong           = "Title must not be more than 30 characters long"
	msgPostContentTooShort    = "Content must be 10 characters"
	msgPostCommentTooLong     = "Content must not be more than 620 characters long"
	msgCommentContentTooShort = "Content must be more than 3 characters long"
	msgCommentContentTooLong  = "Content must not be more than 120 characters long"
)

var PostValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"title": zog.String().Required(zog.Message(msgTitleRequired)).
		Min(4, zog.Message(msgTitleTooShort)).Max(30, zog.Message(msgTitleTooLong)),
	"content": zog.String().Required(zog.Message(msgContentRequired)).
		Min(10, zog.Message(msgPostContentTooShort)).Max(620, zog.Message(msgPostContentTooShort)),
}))

var CreateCommentValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"postId": zog.String().Required(zog.Message(msgPostIdRequired)),
	"content": zog.String().Required(zog.Message(msgContentRequired)).
		Min(10, zog.Message(msgCommentContentTooShort)).Max(120, zog.Message(msgCommentContentTooLong)),
}))

var UpdateCommentValidator = validator.NewStructValidator(zog.Struct(zog.Shape{
	"content": zog.String().Required(zog.Message(msgContentRequired)).
		Min(10, zog.Message(msgCommentContentTooShort)).Max(120, zog.Message(msgCommentContentTooLong)),
}))
