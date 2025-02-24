package repositories

import (
	"CloneVK/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, email, password_hash, avatar_url) VALUES($1,$2,$3,$4) returning userId"
	err := ur.DB.QueryRow(context.Background(), query, user.Username, user.Email, user.PasswordHash, user.AvatarURL).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindUserByID(id int) (*models.User, error) {
	user := models.User{ID: id}
	query := "SELECT username, email, password_hash, avatar_url FROM users WHERE userid=$1"
	err := ur.DB.QueryRow(context.Background(), query, id).Scan(&user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) FindAllUsers() (*[]models.User, error) {
	query := "SELECT userid, username, email, password_hash, avatar_url FROM users"
	rows, err := ur.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
