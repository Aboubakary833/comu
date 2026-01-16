package handlers

import (
	"comu/internal/modules/post/application/posts"

	"github.com/labstack/echo/v4"
)

type postHandlers struct {
	listPostsUC  *posts.ListPostsUC
	readPostUC   *posts.ReadPostUC
	createPostUC *posts.CreatePostUC
	updatePostUC *posts.UpdatePostUC
	deletePostUC *posts.DeletePostUC
}

func newPostHandlers(
	listPostsUC  *posts.ListPostsUC,
	readPostUC   *posts.ReadPostUC,
	createPostUC *posts.CreatePostUC,
	updatePostUC *posts.UpdatePostUC,
	deletePostUC *posts.DeletePostUC,
) *postHandlers {
	return &postHandlers{
		listPostsUC: listPostsUC,
		readPostUC: readPostUC,
		createPostUC: createPostUC,
		updatePostUC: updatePostUC,
		deletePostUC: deletePostUC,
	}
}


func (h *postHandlers) list(ctx echo.Context) error {
	
	return nil
}
