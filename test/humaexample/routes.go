package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pixlcrashr/go-pagetoken"
	"github.com/pixlcrashr/go-pagetoken/encryption"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/repository"
	"gorm.io/gorm"
)

func registerRoutes(api huma.API, db *gorm.DB) {
	k, err := encryption.Rand32ByteKey()
	if err != nil {
		panic(err)
	}

	e, err := encryption.NewAEADEncryptor(k)
	if err != nil {
		panic(err)
	}

	h := &Handler{
		r: &repository.BooksRepository{DB: db},
		rr: pagetoken.NewRequestReader(
			pagetoken.WithEncryptor(e),
		),
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-books-sql",
		Method:      http.MethodGet,
		Path:        "/api/v1/books/sql",
		Summary:     "List books using raw SQL with keyset page-token pagination",
		Tags:        []string{"Books"},
	}, h.ListBooksSQL)
	huma.Register(api, huma.Operation{
		OperationID: "list-books-dao",
		Method:      http.MethodGet,
		Path:        "/api/v1/books/dao",
		Summary:     "List books using DAO with keyset page-token pagination",
		Tags:        []string{"Books"},
	}, h.ListBooksDAO)
}
