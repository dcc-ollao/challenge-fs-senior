package models

type MinimalUser struct {
	ID    string `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
}
