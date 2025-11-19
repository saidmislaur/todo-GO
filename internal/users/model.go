package users

import "database/sql"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserManager struct {
	DB       *sql.DB
	Sessions map[string]int
}
