package users

import (
	"encoding/json"
	"os"
)

func (um *UserManager) SaveUsers() error {
	var users []User
	for _, user := range um.Users {
		users = append(users, user)
	}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(um.FilePath, data, 0644)
}

func (um *UserManager) LoadUsers() error {
	data, err := os.ReadFile(um.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			um.Users = make(map[int]User)
			return nil
		}
		return err
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	um.Users = make(map[int]User)
	for _, task := range users {
		um.Users[task.ID] = task
	}

	return nil
}
