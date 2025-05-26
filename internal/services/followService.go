package services

import (
	followdto "CloneVK/internal/dto/follows"
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type followService struct {
	FollowRepository repositories.IFollowRepositories
}

func NewFollowService(followRepository repositories.IFollowRepositories) IFollowService {
	return &followService{followRepository}
}

func (s *followService) CreateFollow(dto followdto.FollowDTO) error {
	followerID := dto.FollowerID
	followedID := dto.FollowedID
	return s.FollowRepository.CreateFollow(followerID, followedID)
}

func (s *followService) GetAllFollows() ([]models.Follow, error) {
	return s.FollowRepository.GetAllFollows()
}

func (s *followService) GetAllUserFollowers(followedID int) ([]models.Follow, error) {
	return s.FollowRepository.GetAllUserFollowers(followedID)
}

func (s *followService) GetAllUserFollows(followerID int) ([]models.Follow, error) {
	return s.FollowRepository.GetAllUserFollows(followerID)
}

func (s *followService) DeleteFollow(dto followdto.FollowDTO) error {
	followerID := dto.FollowerID
	followedID := dto.FollowedID
	return s.FollowRepository.DeleteFollow(followerID, followedID)
}
