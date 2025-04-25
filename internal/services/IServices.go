package services

import (
	followdto "CloneVK/internal/dto/follows"
	postdto "CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
	"time"
)

type IUserService interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	Register(username, email, password string) error

	Login(email, password string) (*models.User, error)

	SaveRefreshToken(userID int, token string, expiresAt time.Time) error
}

type IPostService interface {
	CreatePost(dto *postdto.CreatePostDTO) (int, error)

	FindPostByID(id int) (*models.Post, error)

	GetAllPostsByUser(userId int) ([]models.Post, error)

	DeletePost(id int) error

	UpdatePost(newPost *postdto.UpdatePostDTO) error
}

type IFeedService interface {
	GetGlobalFeed(limit int, offset int) ([]models.Post, error)
	GetPersonalFeed(userid int, limit int, offset int) ([]models.Post, error)
}

type IFollowService interface {
	CreateFollow(dto followdto.FollowDTO) error

	GetAllFollows() ([]models.Follow, error)

	GetAllUserFollows(followerID int) ([]models.Follow, error)

	GetAllUserFollowers(followedID int) ([]models.Follow, error)

	DeleteFollow(dto followdto.FollowDTO) error
}
