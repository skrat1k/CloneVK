package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type feedService struct {
	FeedRepository repositories.IFeedRepositories
}

// CreateFollow implements IFollowService.
func (s *feedService) CreateFollow(followerID int, followedID int) error {
	panic("unimplemented")
}

// DeleteFollow implements IFollowService.
func (s *feedService) DeleteFollow(followerID int, followedID int) error {
	panic("unimplemented")
}

// GetAllFollows implements IFollowService.
func (s *feedService) GetAllFollows() ([]models.Follow, error) {
	panic("unimplemented")
}

// GetAllUserFollowers implements IFollowService.
func (s *feedService) GetAllUserFollowers(followedID int) ([]models.Follow, error) {
	panic("unimplemented")
}

// GetAllUserFollows implements IFollowService.
func (s *feedService) GetAllUserFollows(followerID int) ([]models.Follow, error) {
	panic("unimplemented")
}

func NewFeedService(feedRepository repositories.IFeedRepositories) IFeedService {
	return &feedService{feedRepository}
}

func (s *feedService) GetGlobalFeed(limit int, offset int) ([]models.Post, error) {
	return s.FeedRepository.GetGlobalFeed(limit, offset)
}
