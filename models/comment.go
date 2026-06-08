package models

type Comments struct {
	ID        int
	PostID    int
	AuthorID  int
	Author    string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
}
