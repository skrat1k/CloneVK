package models

type Post struct {
	ID      int    `json:"postID"`
	UserID  int    `json:"userID"`
	Content string `json:"content"`
	ImgURL  string `json:"imageURL"`
}
