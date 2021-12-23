package store

import (
	"context"

	"github.com/henripqt/lab/pagination/pkg/models"
)

// Repository is an interface that defines the methods that a store must implement
type Repository interface {
	GetBlogPosts(ctx context.Context, paginationReq models.PaginationRequest) (*models.PagingResponse, error)
	Close() error
}

type repository struct {
	repository Repository
}

func NewReposoitory(r Repository) Repository {
	return &repository{
		repository: r,
	}
}

var _ Repository = (*repository)(nil)

func (r *repository) GetBlogPosts(ctx context.Context, paginationReq models.PaginationRequest) (*models.PagingResponse, error) {
	return r.repository.GetBlogPosts(ctx, paginationReq)
}

func (r *repository) Close() error {
	return r.repository.Close()
}
