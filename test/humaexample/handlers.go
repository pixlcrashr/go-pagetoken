package main

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pixlcrashr/go-pagetoken"
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

	order := []repository.DefaultOrderEntry{}
	if req.OrderBy != "" {
		order = lo.Map(strings.Split(req.OrderBy, ","), func(s string, _ int) repository.DefaultOrderEntry {
			ps := strings.Split(s, " ")
			if len(ps) == 1 {
				ps = append(ps, "DESC")
			}
			o, err := pagetoken.ParseOrder(ps[1])
			if err != nil {
				panic(err)
			}
			return repository.DefaultOrderEntry{
				Column: ps[0],
				Order:  o,
			}
		})
	}

	ms, nextPayload, err := h.r.ListByKeyset(
		ctx,
		repository.ListFilter{},
		req.PageSize,
		order,
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
