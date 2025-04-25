package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type feedService struct {
	FeedRepository repositories.IFeedRepositories
}

func NewFeedService(feedRepository repositories.IFeedRepositories) IFeedService {
	return &feedService{feedRepository}
}

func (s *feedService) GetGlobalFeed(limit int, offset int) ([]models.Post, error) {
	return s.FeedRepository.GetGlobalFeed(limit, offset)
}

func (s *feedService) GetPersonalFeed(userid int, limit int, offset int) ([]models.Post, error) {
	return s.FeedRepository.GetPersonalFeed(userid, limit, offset)
}
