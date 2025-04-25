package repositories

import (
	"CloneVK/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type followRepository struct {
	DB *pgx.Conn
}

func NewFollowRepository(database *pgx.Conn) IFollowRepositories {
	return &followRepository{database}
}

func (r *followRepository) CreateFollow(followerID int, followedID int) error {
	query := "INSERT INTO follows (follower_id, followed_id) VALUES($1, $2)"
	_, err := r.DB.Exec(context.Background(), query, followerID, followedID)
	if err != nil {
		return err
	}
	return nil
}

// TODO: мб в будущем добавить лимиты и оффсеты
func (r *followRepository) GetAllFollows() ([]models.Follow, error) {
	query := "SELECT id, follower_id, followed_id FROM follows"
	rows, err := r.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []models.Follow

	for rows.Next() {
		follow := models.Follow{}
		err := rows.Scan(&follow.ID, &follow.FollowerID, &follow.FollowedID)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

func (r *followRepository) GetAllUserFollows(followerID int) ([]models.Follow, error) {
	query := "SELECT id, followed_id FROM follows WHERE follower_id = $1"
	rows, err := r.DB.Query(context.Background(), query, followerID)
	if err != nil {
		return nil, err
	}

	var follows []models.Follow

	for rows.Next() {
		follow := models.Follow{FollowerID: followerID}
		err := rows.Scan(&follow.ID, &follow.FollowedID)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

func (r *followRepository) GetAllUserFollowers(followedID int) ([]models.Follow, error) {
	query := "SELECT id, followed_id FROM follows WHERE followed_id = $1"
	rows, err := r.DB.Query(context.Background(), query, followedID)
	if err != nil {
		return nil, err
	}

	var follows []models.Follow

	for rows.Next() {
		follow := models.Follow{FollowerID: followedID}
		err := rows.Scan(&follow.ID, &follow.FollowedID)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

func (r *followRepository) DeleteFollow(followerID int, followedID int) error {
	query := "DELETE FROM follows WHERE follower_ID = $1 AND followed_ID = $2"
	_, err := r.DB.Exec(context.Background(), query, followerID, followedID)
	if err != nil {
		return err
	}
	return nil
}
