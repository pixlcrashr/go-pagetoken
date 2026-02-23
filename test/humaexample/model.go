package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/pixlcrashr/go-pagetoken/checksum"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
)

// Book is the API representation of a Book.
type Book struct {
	ID          uuid.UUID `json:"id" doc:"Book UUID"`
	DisplayName string    `json:"display_name" doc:"Human-readable book name"`
	CreatedAt   time.Time `json:"created_at" doc:"Creation timestamp"`
	UpdatedAt   time.Time `json:"updated_at" doc:"Last modification timestamp"`
}

func (b *Book) fromModel(m *model.Book) {
	b.ID = m.ID
	b.DisplayName = m.DisplayName
	b.CreatedAt = m.CreatedAt
	b.UpdatedAt = m.UpdatedAt
}

// ListBooksRequest holds query parameters for the list-books endpoint.
type ListBooksRequest struct {
	DisplayName string `query:"display_name" doc:"Case-sensitive prefix filter on display_name" maxLength:"200"`
	ID          string `query:"id" doc:"Filter by exact book UUID" maxLength:"36"`
	OrderBy     string `query:"order_by" doc:"Sort expression, e.g. 'display_name desc'. Available fields: id, display_name, created_at, updated_at"`
	PageSize    int    `query:"page_size" doc:"Books per page (max 100)" minimum:"1" maximum:"100" default:"20"`
	PageToken   string `query:"page_token" doc:"Opaque continuation token from the previous response; omit for the first page" maxLength:"512"`
}

func (r *ListBooksRequest) GetPageToken() string { return r.PageToken }
func (r *ListBooksRequest) GetChecksumFields() []checksum.BuilderOpt {
	return []checksum.BuilderOpt{
		checksum.Field("display_name", r.DisplayName),
		checksum.Field("id", r.ID),
		checksum.Field("order_by", r.OrderBy),
	}
}

// ListBooksResponse is the API response for the list-books endpoint.
type ListBooksResponse struct {
	Body struct {
		Books         []Book `json:"books"`
		NextPageToken string `json:"next_page_token" doc:"Pass as page_token on the next request; empty on the last page" maxLength:"512"`
	}
}
