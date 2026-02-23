package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pixlcrashr/go-pagetoken"
	ptGorm "github.com/pixlcrashr/go-pagetoken/database/gorm"
	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
	"gorm.io/gorm"
)

type BooksRepository struct {
	DB *gorm.DB
}

type ListFilter struct {
	DisplayNameEq   *string
	DisplayNameLike *string
	IDEq            *string
}

type DefaultOrderEntry struct {
	Column string
	Order  pagetoken.Order
}

func orderToSQL(o pagetoken.Order) string {
	if o == pagetoken.OrderAsc {
		return "ASC"
	}

	return "DESC"
}

func (r *BooksRepository) ListByKeyset(
	ctx context.Context,
	filter ListFilter,
	pageSize int,
	order []DefaultOrderEntry,
	keyset *pagetoken.KeysetPayload,
) (ms []*model.Book, next *pagetoken.KeysetPayload, err error) {
	q := r.DB.Model(&model.Book{})

	if filter.DisplayNameEq != nil {
		q = q.Where("display_name = ?", *filter.DisplayNameEq)
	}

	if filter.DisplayNameLike != nil {
		q = q.Where("display_name LIKE ?", *filter.DisplayNameLike)
	}

	if filter.IDEq != nil {
		q = q.Where("id = ?", *filter.IDEq)
	}

	q, err = ptGorm.KeysetWhereOrderLimit(q, keyset, func(column string, payload *pagetoken.KeysetPayload) (any, error) {
		switch column {
		case "id":
			v, _, err := payload.String(column)
			if err != nil {
				return nil, err
			}

			id, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}

			return id, nil
		case "display_name":
			v, _, err := payload.String(column)
			if err != nil {
				return nil, err
			}
			return v, nil
		case "created_at":
			v, _, err := payload.Time(column)
			if err != nil {
				return nil, err
			}
			return v, nil
		default:
			return nil, errors.New("unknown column: " + column)
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply keyset: %w", err)
	}

	if keyset == nil || len(keyset.Values()) == 0 {
		if len(order) == 0 {
			q = q.Order("created_at DESC")
		} else {
			for _, o := range order {
				q = q.Order(fmt.Sprintf("%s %s", o.Column, orderToSQL(o.Order)))
			}
		}
	}

	q = q.Limit(pageSize + 1)

	ms = []*model.Book{}
	if err := q.Find(&ms).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to query books: %w", err)
	}

	// TODO: handle this part more DX friendly
	if len(ms) > pageSize {
		nextItem := ms[pageSize-1]

		if keyset == nil || len(keyset.Values()) == 0 {
			if len(order) == 0 {
				keyset = pagetoken.NewKeysetPayloadBuilder().
					AddTime("created_at", time.Time{}, pagetoken.OrderDesc).
					Build()
			} else {
				b := pagetoken.NewKeysetPayloadBuilder()
				// empty default setter
				func(m *model.Book, keysetBuilder *pagetoken.KeysetPayloadBuilder, values []DefaultOrderEntry) {
					for _, v := range values {
						switch v.Column {
						case "id":
							keysetBuilder.AddString(v.Column, "", v.Order)
						case "display_name":
							keysetBuilder.AddString(v.Column, "", v.Order)
						case "created_at":
							keysetBuilder.AddTime(v.Column, time.Time{}, v.Order)
						}
					}
				}(nextItem, b, order)
				keyset = b.Build()
			}
		}

		keysetBuilder := pagetoken.NewKeysetPayloadBuilder()

		func(m *model.Book, keysetBuilder *pagetoken.KeysetPayloadBuilder, values []pagetoken.KeysetValue) {
			for _, v := range values {
				switch v.Path {
				case "id":
					keysetBuilder.AddString(v.Path, m.ID.String(), v.Order)
				case "display_name":
					keysetBuilder.AddString(v.Path, m.DisplayName, v.Order)
				case "created_at":
					keysetBuilder.AddTime(v.Path, m.CreatedAt, v.Order)
				}
			}
		}(nextItem, keysetBuilder, keyset.Values())

		next = keysetBuilder.Build()
	}

	return ms[:minPageSize(len(ms), pageSize)], next, nil
}

func minPageSize(a, b int) int {
	if a < b {
		return a
	}
	return b
}
