package store

import (
	"context"

	"github.com/henripqt/lab/pagination/pkg/models"
)

// Repository is an interface that defines the methods that a store must implement
type Repository interface {
	GetBlogPosts(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogPost], error)
	GetBlogCategories(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogCategory], error)
	Close() error
}

// repository is the concrete implementation of the Repository interface
type repository struct {
	repository Repository
}

// NewRepository returns a new instance of the Repository interface
func NewReposoitory(r Repository) Repository {
	return &repository{
		repository: r,
	}
}

var _ Repository = (*repository)(nil)

func (r *repository) GetBlogPosts(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogPost], error) {
	return r.repository.GetBlogPosts(ctx, paginationReq)
}

func (r *repository) GetBlogCategories(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogCategory], error) {
	return r.repository.GetBlogCategories(ctx, paginationReq)
}

func (r *repository) Close() error {
	return r.repository.Close()
}
