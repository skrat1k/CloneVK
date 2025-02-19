package models

type Comment struct {
	ID      int
	PostID  string
	UserID  string
	Content string
}
