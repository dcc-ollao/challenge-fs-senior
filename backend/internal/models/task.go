package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `db:"id"`
	ProjectID   uuid.UUID  `db:"project_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	AssigneeID  *uuid.UUID `db:"assignee_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
