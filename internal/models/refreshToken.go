package models

type RefreshToken struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}
