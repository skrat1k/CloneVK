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

func (r *postRepository) FindPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT postid, userid, post_content, image_url FROM posts WHERE postid = $1"
	err := r.DB.QueryRow(context.Background(), query, id).Scan(&post.ID, &post.UserID, &post.Content, &post.ImgURL)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) GetAllPostsByUser(userId int) ([]models.Post, error) {
	query := "SELECT postid, post_content, image_url FROM posts WHERE userid = $1"
	rows, err := r.DB.Query(context.Background(), query, userId)
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

	return posts, nil
}

func (r *postRepository) CreatePost(post *models.Post) error {
	query := "INSERT INTO posts (userID, post_content, image_url) VALUES($1,$2,$3) returning postID"
	err := r.DB.QueryRow(context.Background(), query, post.UserID, post.Content, post.ImgURL).Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) DeletePost(id int) error {
	query := "DELETE FROM posts WHERE postid = $1"
	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) UpdatePost(post *models.Post) error {
	query := "UPDATE posts SET post_content = $1, image_url = $2 WHERE postid = $3"
	_, err := r.DB.Exec(context.Background(), query, post.Content, post.ImgURL, post.ID)
	if err != nil {
		return err
	}
	return nil
}
