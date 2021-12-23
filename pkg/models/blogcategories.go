package models

import "time"

// BlogCategory is a struct that represents a blog category
type BlogCategory struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
