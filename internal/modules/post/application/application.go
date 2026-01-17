package application

import (
	"comu/internal/modules/post/application/comments"
	"comu/internal/modules/post/application/posts"
	"comu/internal/modules/post/domain"
)

type UseCases struct {
	ListPostsUC  *posts.ListPostsUC
	ReadPostUC   *posts.ReadPostUC
	CreatePostUC *posts.CreatePostUC
	UpdatePostUC *posts.UpdatePostUC
	DeletePostUC *posts.DeletePostUC

	ListCommentUC   *comments.ListCommentsUC
	CreateCommentUC *comments.CreateCommentUC
	UpdateCommentUC *comments.UpdateCommentUC
	DeleteCommentUC *comments.DeleteCommentUC
}

func InitUseCases(
	postsRepository domain.PostRepository,
	commentRepository domain.CommentRepository,
) UseCases {

	readPostUC := posts.NewReadPostUseCase(postsRepository)
	listPostsUC := posts.NewListPostsUseCase(postsRepository)
	createPostUC := posts.NewCreatePostUseCase(postsRepository)
	updatePostUC := posts.NewUpdatePostUseCase(postsRepository)
	deletePostUC := posts.NewDeletePostUseCase(postsRepository)

	listCommentsUC := comments.NewListCommentsUseCase(commentRepository)
	createCommentUC := comments.NewCreateCommentUseCase(commentRepository)
	updateCommentUC := comments.NewUpdateCommentUseCase(commentRepository)
	deleteCommentUC := comments.NewDeleteCommentUseCase(commentRepository)

	return UseCases{
		ListPostsUC:  listPostsUC,
		ReadPostUC:   readPostUC,
		CreatePostUC: createPostUC,
		UpdatePostUC: updatePostUC,
		DeletePostUC: deletePostUC,

		ListCommentUC:   listCommentsUC,
		CreateCommentUC: createCommentUC,
		UpdateCommentUC: updateCommentUC,
		DeleteCommentUC: deleteCommentUC,
	}
}
