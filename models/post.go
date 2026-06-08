package models

type Post struct {
	ID        int
	Title     string
	AuthorID  int
	Author    string
	Content   string
	Category  string
	ImagePath string
	Likes     int
	Dislikes  int
	CreatedAt string
}
