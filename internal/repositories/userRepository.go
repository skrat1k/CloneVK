package repositories

import (
	"CloneVK/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	DB *pgx.Conn
}

func NewUserRepositories(database *pgx.Conn) IUserRepositories {
	return &userRepository{database}
}

func (ur *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, email, password_hash, avatar_url) VALUES($1,$2,$3,$4) returning userId"
	err := ur.DB.QueryRow(context.Background(), query, user.Username, user.Email, user.PasswordHash, user.AvatarURL).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindUserByID(id int) (*models.User, error) {
	user := models.User{ID: id}
	query := "SELECT username, email, password_hash, avatar_url FROM users WHERE userid=$1"
	err := ur.DB.QueryRow(context.Background(), query, id).Scan(&user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindAllUsers() (*[]models.User, error) {
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

func (ur *userRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT userid, username, email, password_hash, avatar_url FROM users WHERE email=$1"
	err := ur.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := ur.DB.Exec(context.Background(), query, userID, token, expiresAt)
	return err
}

func (ur *userRepository) DeleteRefreshTokensByUserID(userID int) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = $1"
	_, err := ur.DB.Exec(context.Background(), query, userID)
	return err
}
