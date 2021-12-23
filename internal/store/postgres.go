package store

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/henripqt/lab/pagination/pkg/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type PGPaginator[T any] []T

type PGRepository struct {
	sq squirrel.StatementBuilderType
	db *sqlx.DB
}

var _ Repository = (*PGRepository)(nil)

func NewPGRepository(userName, password, dbName string) Repository {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", userName, password, dbName))
	if err != nil {
		log.Fatalln(err)
	}

	return &PGRepository{
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db: db,
	}
}

// GetBLogPosts returns all blog posts
func (r *PGRepository) GetBlogPosts(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogPost], error) {
	query, queryArgs, err := r.sq.Select("*").From("blog_posts").ToSql()
	if err != nil {
		return nil, err
	}

	countQuery, countQueryArgs, err := r.sq.Select("count(*)").From("blog_posts").ToSql()
	if err != nil {
		return nil, err
	}

	pRes, err := PGPaginator[models.BlogPost]{}.paginate(
		r.db,
		ctx,
		query,
		queryArgs,
		countQuery,
		countQueryArgs,
		paginationReq,
	)

	if err != nil {
		return nil, err
	}

	return pRes, nil
}

// GetBLogPosts returns all blog posts
func (r *PGRepository) GetBlogCategories(ctx context.Context, paginationReq models.PaginationRequest) (*models.PaginationResponse[[]models.BlogCategory], error) {
	query, queryArgs, err := r.sq.Select("*").From("blog_categories").ToSql()
	if err != nil {
		return nil, err
	}

	countQuery, countQueryArgs, err := r.sq.Select("count(*)").From("blog_categories").ToSql()
	if err != nil {
		return nil, err
	}

	pRes, err := PGPaginator[models.BlogCategory]{}.paginate(
		r.db,
		ctx,
		query,
		queryArgs,
		countQuery,
		countQueryArgs,
		paginationReq,
	)

	if err != nil {
		return nil, err
	}

	return pRes, nil
}

// Close allows for closing the database connection
func (r *PGRepository) Close() error {
	return r.db.Close()
}

// paginate is a helper function for fetching paginated ressources
func (r PGPaginator[T]) paginate(
	db *sqlx.DB,
	ctx context.Context,
	query string,
	queryArgs []interface{},
	countQuery string,
	countQueryArgs []interface{},
	paginationReq models.PaginationRequest,
) (*models.PaginationResponse[[]T], error) {
	paginationRes := models.PaginationResponse[[]T]{
		Page:    paginationReq.Page,
		PerPage: paginationReq.PerPage,
	}

	g, _ := errgroup.WithContext(ctx)

	// Retrieve the total number of items
	g.Go(func() error {
		return db.GetContext(
			ctx,
			&paginationRes.TotalItems,
			countQuery,
			countQueryArgs...,
		)
	})

	// Retrieve the items
	var rows *sqlx.Rows
	g.Go(func() error {
		var err error
		rows, err = db.QueryxContext(
			ctx,
			r.decoratePaginatedQuery(query, paginationReq),
			queryArgs...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		var item T
		err := rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	paginationRes.Items = items
	paginationRes.TotalPage = r.getTotalPage(int(paginationRes.TotalItems), paginationRes.PerPage)
	paginationRes.PrevPage = r.getPrevPage(paginationRes.Page)
	paginationRes.NextPage = r.getNextPage(paginationRes.Page, paginationRes.TotalPage)
	return &paginationRes, nil
}

// getTotalPage is a helper function for getting the total number of pages
func (r PGPaginator[T]) getTotalPage(totalItems, perPage int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

// getPrevPage is a helper function for getting the previous page
func (r PGPaginator[T]) getPrevPage(currentPage int) int {
	if currentPage >= 2 {
		return currentPage - 1
	}
	return currentPage
}

// getNextPage is a helper function for getting the next page
func (r PGPaginator[T]) getNextPage(currentPage, totalPage int) int {
	if currentPage >= totalPage {
		return currentPage
	}
	return currentPage + 1
}

func (r PGPaginator[T]) decoratePaginatedQuery(query string, pReq models.PaginationRequest) string {
	q := strings.Builder{}
	q.WriteString(query)

	if len(pReq.OrderBy) > 0 {
		// ORDER BY instructions
		q.WriteRune(' ')
		q.WriteString("ORDER BY")

		for i, orderBy := range pReq.OrderBy {
			if i > 0 {
				q.WriteRune(',')
			}
			q.WriteRune(' ')
			q.WriteString(orderBy)
		}

		q.WriteRune(' ')
		if len(pReq.OrderDir) == 0 {
			q.WriteString("DESC")
		} else {
			q.WriteString(pReq.OrderDir)
		}
	}

	// LIMIT instruction
	q.WriteRune(' ')
	q.WriteString("LIMIT")
	q.WriteRune(' ')
	q.WriteString(strconv.Itoa(pReq.PerPage))

	// OFFSET instruction
	q.WriteRune(' ')
	q.WriteString("OFFSET")
	q.WriteRune(' ')
	q.WriteString(strconv.Itoa(pReq.PerPage * (pReq.Page - 1)))

	return q.String()
}
