package main

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pixlcrashr/go-pagetoken"
	"github.com/pixlcrashr/go-pagetoken/order"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/repository"
	"github.com/samber/lo"
)

var (
	ErrInvalidPageToken  = huma.Error400BadRequest("invalid page_token")
	ErrInvalidOrderBy    = huma.Error400BadRequest("invalid order_by")
	ErrFailedToListBooks = huma.Error500InternalServerError("failed to list books")
)

// Handler holds the dependencies for the books API.
type Handler struct {
	r  *repository.BooksRepository
	rr *pagetoken.RequestReader
}

// ListBooks handles GET /api/v1/books.
// TODO: implement keyset-pagination using pagetoken.KeysetTokenFromRequest.
func (h *Handler) ListBooksSQL(_ context.Context, _ *ListBooksRequest) (*ListBooksResponse, error) {
	return nil, ErrFailedToListBooks
}

func (h *Handler) ListBooksDAO(ctx context.Context, req *ListBooksRequest) (*ListBooksResponse, error) {
	t, err := h.rr.Read(req)
	if err != nil {
		return nil, ErrInvalidPageToken
	}

	oFs := order.Fields{}
	if req.OrderBy != "" {
		if err := oFs.UnmarshalString(req.OrderBy); err != nil {
			return nil, ErrInvalidOrderBy
		}
	}

	ms, nextPayload, err := h.r.ListByKeyset(
		ctx,
		repository.ListFilter{},
		req.PageSize,
		oFs,
		t.Payload(),
	)
	if err != nil {
		return nil, ErrFailedToListBooks
	}

	resp := &ListBooksResponse{}
	if nextPayload != nil {
		t = t.Next(pagetoken.WithKeysetPayload(nextPayload))

		tS, err := t.String()
		if err != nil {
			return nil, ErrInvalidPageToken
		}

		resp.Body.NextPageToken = tS
	}

	resp.Body.Books = lo.Map(ms, func(m *model.Book, _ int) Book {
		b := Book{}
		b.fromModel(m)
		return b
	})

	return resp, nil
}
