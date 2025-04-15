package repositories

import (
	"CloneVK/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type postRepository struct {
	DB *pgx.Conn
}

func NewPostRepository(database *pgx.Conn) IPostRepositories {
	return &postRepository{database}
}

func (pr *postRepository) FindPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT postid, userid, post_content, image_url FROM posts"
	err := pr.DB.QueryRow(context.Background(), query).Scan(&post.ID, &post.UserID, &post.Content, &post.ImgURL)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *postRepository) GetAllPostsByUser(userId int) (*[]models.Post, error) {
	query := "SELECT postid, post_content, image_url FROM posts WHERE userid = $1"
	rows, err := pr.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		post := models.Post{UserID: userId}
		err := rows.Scan(&post.ID, &post.Content, &post.ImgURL)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (pr *postRepository) CreatePost(post *models.Post) error {
	query := "INSERT INTO posts (userID, post_content, image_url) VALUES($1,$2,$3) returning postID"
	err := pr.DB.QueryRow(context.Background(), query, post.UserID, post.Content, post.ImgURL).Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}
