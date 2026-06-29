package domain

import "time"

type AuthUser struct {
	ID           string
	ProfileID    string
	ProfileType  string
	Email        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type UserInfo struct {
	ID          string
	Email       string
	Name        string
	ProfileType string
	Profile     interface{}
}
