package models

type User struct {
	ID           int    `json:"userID"`
	Username     string `json:"userName"`
	Email        string `json:"userEmail"`
	PasswordHash string `json:"userPassword"`
	AvatarURL    string `json:"userAvatar"`
}
