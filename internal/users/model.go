package users

import "bufio"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserManager struct {
	Users    map[int]User
	Reader   *bufio.Reader
	FilePath string
	Sessions map[string]int
}
