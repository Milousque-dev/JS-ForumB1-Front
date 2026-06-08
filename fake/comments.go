package fake

import (
	"forum/database"
	"forum/models"
)

func GetCommentByPostID(postID int) []models.Comments {
	return database.GetCommentsByPostID(postID)
}

func GetCommentByID(commentID int) (models.Comments, bool) {
	return database.GetCommentByID(commentID)
}
