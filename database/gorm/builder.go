package gorm

import (
	"fmt"
	"strings"

	"github.com/pixlcrashr/go-pagetoken"
	"gorm.io/gorm"
)

type KeysetWhereOrderLimitValueFn func(column string, payload *pagetoken.KeysetPayload) (any, error)

func orderToSQL(o pagetoken.Order) string {
	if o == pagetoken.OrderAsc {
		return "ASC"
	}

	return "DESC"
}

func KeysetWhereOrderLimit(
	db *gorm.DB,
	keyset *pagetoken.KeysetPayload,
	valueFn KeysetWhereOrderLimitValueFn,
) (*gorm.DB, error) {
	if keyset == nil {
		return db, nil
	}

	vs := keyset.Values()
	if len(vs) == 0 {
		return db, nil
	}

	args := []any{}
	orExprs := []string{}
	orderExprs := []string{}

	for i := 0; i < len(vs); i++ {
		andExprs := []string{}
		for j := 0; j < i; j++ {
			v := vs[j]
			andExprs = append(andExprs, fmt.Sprintf("%s = ?", v.Path))
			aV, err := valueFn(v.Path, keyset)
			if err != nil {
				return nil, err
			}
			args = append(args, aV)
		}

		v := vs[i]
		if v.Order == pagetoken.OrderDesc {
			andExprs = append(andExprs, fmt.Sprintf("%s < ?", v.Path))
			aV, err := valueFn(v.Path, keyset)
			if err != nil {
				return nil, err
			}
			args = append(args, aV)
		} else {
			andExprs = append(andExprs, fmt.Sprintf("%s > ?", v.Path))
			aV, err := valueFn(v.Path, keyset)
			if err != nil {
				return nil, err
			}
			args = append(args, aV)
		}

		orExprs = append(orExprs, "("+strings.Join(andExprs, " AND ")+")")
		orderExprs = append(orderExprs, fmt.Sprintf("%s %s", v.Path, orderToSQL(v.Order)))
	}

	return db.Where(
		"("+strings.Join(orExprs, " OR ")+")",
		args...,
	).Order(
		strings.Join(orderExprs, ", "),
	), nil
}
