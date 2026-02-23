package gorm

import (
	"github.com/pixlcrashr/go-pagetoken"
	"gorm.io/gorm/clause"
)

type KeysetValueOp int

const (
	KeysetValueOpEq KeysetValueOp = iota
	KeysetValueOpLt
	KeysetValueOpGt
)

type KeysetValueFn func(field string, value string, op KeysetValueOp) (clause.Expression, error)

func KeysetValuesExpr(fields []pagetoken.KeysetValue, fieldFn KeysetValueFn) (clause.Expression, error) {
	or := []clause.Expression{}

	for i := 0; i < len(fields); i++ {
		and := []clause.Expression{}

		for j := 0; j < i; j++ {
			f := fields[j]
			expr, err := fieldFn(f.Path, f.Value, KeysetValueOpEq)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		}

		f := fields[i]
		if f.Order == pagetoken.OrderDesc {
			expr, err := fieldFn(f.Path, f.Value, KeysetValueOpLt)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		} else {
			expr, err := fieldFn(f.Path, f.Value, KeysetValueOpGt)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		}

		or = append(or, clause.And(and...))
	}

	return clause.Or(or...), nil
}
