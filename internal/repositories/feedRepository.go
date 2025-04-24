package repositories

import (
	"CloneVK/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type feedRepository struct {
	DB *pgx.Conn
}

func NewFeedRepository(database *pgx.Conn) IFeedRepositories {
	return &feedRepository{database}
}

func (r *feedRepository) GetGlobalFeed(limit int, offset int) ([]models.Post, error) {
	query := "SELECT postid, userid, post_content, image_url FROM posts LIMIT $1 OFFSET $2"
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0, limit)

	for rows.Next() {
		post := models.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImgURL)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
