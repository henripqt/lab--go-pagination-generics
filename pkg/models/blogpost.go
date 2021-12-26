package models

import (
	"fmt"
	"time"
)

// BlogPost is a struct that represents a blog post
type BlogPost struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Body      string    `json:"body" db:"body"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (b BlogPost) BlogPostMethod() string {
	return fmt.Sprintf("I'm a blog post ! Title: %v", b.Title)
}
