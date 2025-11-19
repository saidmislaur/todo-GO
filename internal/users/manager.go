package users

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func NewManager(db *sql.DB) *UserManager {
	return &UserManager{
		DB:       db,
		Sessions: make(map[string]int),
	}
}

func (um *UserManager) Register(username, password string) error {
	var existingID int
	err := um.DB.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("пользователь %q уже существует", username)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("ошибка проверки пользователя: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хэширования пароля: %v", err)
	}

	_, err = um.DB.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		username, string(hash),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	return nil
}

func (um *UserManager) Login(username, password string) (string, error) {
	var id int
	var storedHash string

	err := um.DB.QueryRow(
		"SELECT id, password_hash FROM users WHERE username=$1",
		username,
	).Scan(&id, &storedHash)

	if err == sql.ErrNoRows {
		return "", errors.New("неверный логин или пароль")
	}
	if err != nil {
		return "", fmt.Errorf("ошибка чтения пользователя: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return "", errors.New("неверный логин или пароль")
	}

	token := fmt.Sprintf("token-%d-%s", id, username)

	um.Sessions[token] = id

	return token, nil
}

func (um *UserManager) GetUserIDByToken(token string) (int, error) {
	userID, ok := um.Sessions[token]
	if !ok {
		return 0, errors.New("неверный или просроченный токен")
	}
	return userID, nil
}
