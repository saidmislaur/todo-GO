package users

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func NewManager(filePath string) *UserManager {
	return &UserManager{
		Reader:   bufio.NewReader(os.Stdin),
		FilePath: filePath,
		Sessions: make(map[string]int),
	}
}

func (um *UserManager) Register(username, password string) error {
	for _, u := range um.Users {
		if u.Username == username {
			return fmt.Errorf("пользователь %q уже существует", username)
		}
	}

	newID := 1
	if len(um.Users) > 0 {
		for id := range um.Users {
			if id >= newID {
				newID = id + 1
			}
		}
	}

	newUser := User{
		ID:       newID,
		Username: username,
		Password: password,
	}

	um.Users[newID] = newUser

	if err := um.SaveUsers(); err != nil {
		return fmt.Errorf("ошибка сохранения пользователя %v", err)
	}

	return nil
}

func (um *UserManager) Login(username, password string) (string, error) {
	for _, user := range um.Users {
		if user.Username == username && user.Password == password {
			token := fmt.Sprintf("%s-%d", username, user.ID)
			um.Sessions[token] = user.ID
			return token, nil
		}
	}

	return "", fmt.Errorf("неверный логин или пароль")
}

func (um *UserManager) GetUserIDByToken(token string) (int, error) {
	id, ok := um.Sessions[token]
	if !ok {
		return 0, errors.New("неверный или просроченный токен")
	}
	return id, nil
}
