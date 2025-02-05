package models


type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}