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
