package gorm

import (
	"github.com/pixlcrashr/go-pagetoken"
	"gorm.io/gorm/clause"
)

type KeysetFieldOp int

const (
	KeysetFieldOpEq KeysetFieldOp = iota
	KeysetFieldOpLt
	KeysetFieldOpGt
)

type KeysetFieldFn func(field string, value string, op KeysetFieldOp) (clause.Expression, error)

func KeysetFieldsExpr(fields []pagetoken.KeysetField, fieldFn KeysetFieldFn) (clause.Expression, error) {
	or := []clause.Expression{}

	for i := 0; i < len(fields); i++ {
		and := []clause.Expression{}

		for j := 0; j < i; j++ {
			f := fields[j]
			expr, err := fieldFn(f.Path, f.Value, KeysetFieldOpEq)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		}

		f := fields[i]
		if f.Order == pagetoken.OrderDesc {
			expr, err := fieldFn(f.Path, f.Value, KeysetFieldOpLt)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		} else {
			expr, err := fieldFn(f.Path, f.Value, KeysetFieldOpGt)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		}

		or = append(or, clause.And(and...))
	}

	return clause.Or(or...), nil
}
