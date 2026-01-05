package models

import "time"

type Project struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OwnerID   string    `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
