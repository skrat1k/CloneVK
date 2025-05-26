package dto

type FollowDTO struct {
	FollowerID int `json:"followerID" validate:"required"`
	FollowedID int `json:"followedID" validate:"required"`
}
