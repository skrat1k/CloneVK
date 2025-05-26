package repositories

import (
	"CloneVK/internal/models"
	"time"
)

type IUserRepositories interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	FindUserByEmail(email string) (*models.User, error)

	SaveRefreshToken(userID int, token string, expiresAt time.Time) error

	//GetRefreshToken(userID int) (string, error)

	// DeleteRefreshTokensByUserID(userID int) error
}

type IPostRepositories interface {
	FindPostByID(id int) (*models.Post, error)

	GetAllPostsByUser(userId int, limit int, offset int) ([]models.Post, error)

	CreatePost(post *models.Post) error

	DeletePost(id int) error

	UpdatePost(post *models.Post) error
}

type IFeedRepositories interface {
	GetGlobalFeed(limit int, offset int) ([]models.Post, error)
	GetPersonalFeed(userid int, limit int, offset int) ([]models.Post, error)
}

type IFollowRepositories interface {
	CreateFollow(followerID int, followedID int) error

	GetAllFollows() ([]models.Follow, error)

	GetAllUserFollows(followerID int) ([]models.Follow, error)

	GetAllUserFollowers(followedID int) ([]models.Follow, error)

	DeleteFollow(followerID int, followedID int) error
}
