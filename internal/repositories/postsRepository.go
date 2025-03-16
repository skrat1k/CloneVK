package repositories

import (
	"CloneVK/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type PostRepository struct {
	DB *pgx.Conn
}

func (pr *PostRepository) FindPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT postid, userid, post_content, image_url FROM posts"
	err := pr.DB.QueryRow(context.Background(), query).Scan(&post.ID, &post.UserID, &post.Content, &post.ImgURL)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *PostRepository) CreatePost(post *models.Post) error {
	query := "INSERT INTO posts (userID, post_content, image_url) VALUES($1,$2,$3) returning postID"
	err := pr.DB.QueryRow(context.Background(), query, post.UserID, post.Content, post.ImgURL).Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}

/*
func (pr *PostRepository) FindAllPostByUser(id int) (*[]models.Post, error) {
	query := "SELECT postid, post_content, image_url FROM posts WHERE userid = &1"
	rows, err := pr.DB.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := &[]models.Post{}

	for rows.Next() {
		var post models.Post
		err := rows.Scan()
		posts := append(*posts)
	}
}
*/
