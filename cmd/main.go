package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/henripqt/lab/pagination/internal/store"
	"github.com/henripqt/lab/pagination/pkg/models"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	api := &API{
		repository: store.NewPGRepository(
			"postgres",
			"mysecretpassword",
			"mydb",
		),
	}

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: api.handler(),
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}

type API struct {
	repository store.Repository
}

func (a *API) handler() http.Handler {
	router := chi.NewMux()
	router.Get("/blog/categories", a.blogCategoriesHandler)
	router.Get("/blog/posts", a.blogPostsHandler)
	return router
}

func (a *API) blogCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	paginationReq := a.parsePaginationReq(r)

	paginationResponse, err := a.repository.GetBlogCategories(r.Context(), paginationReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, post := range paginationResponse.Items {
		fmt.Println(post.BlogCategoryMethod())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(paginationResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *API) blogPostsHandler(w http.ResponseWriter, r *http.Request) {
	paginationReq := a.parsePaginationReq(r)

	paginationResponse, err := a.repository.GetBlogPosts(r.Context(), paginationReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, post := range paginationResponse.Items {
		fmt.Println(post.BlogPostMethod())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(paginationResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *API) parsePaginationReq(r *http.Request) models.PaginationRequest {
	sPage := r.URL.Query().Get("page")
	page, err := strconv.Atoi(sPage)
	if err != nil {
		page = 1
	}

	sPerPage := r.URL.Query().Get("per_page")
	perPage, err := strconv.Atoi(sPerPage)
	if err != nil {
		perPage = 10
	}

	orderBy := make([]string, 0)
	for _, orderByParam := range r.URL.Query()["order_by"] {
		orderBy = append(orderBy, orderByParam)
	}

	return models.PaginationRequest{
		Page:     page,
		PerPage:  perPage,
		OrderBy:  orderBy,
		OrderDir: r.URL.Query().Get("order_dir"),
	}
}
