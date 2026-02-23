package gorm

import (
	"github.com/pixlcrashr/go-pagetoken"
	"gorm.io/gen/field"
)

type KeysetFieldOp int

const (
	KeysetFieldOpEq KeysetFieldOp = iota
	KeysetFieldOpLt
	KeysetFieldOpGt
)

type KeysetFieldFn func(field string, value string, op KeysetFieldOp) (field.Expr, error)

func KeysetTokenCond(token *pagetoken.KeysetToken, fieldFn KeysetFieldFn) (field.Expr, error) {
	fs := token.Fields()

	or := []field.Expr{}

	for i := 0; i < len(fs)-1; i++ {
		and := []field.Expr{}

		for j := 0; j < i-1; j++ {
			f := fs[j]
			expr, err := fieldFn(f.Path, f.Value, KeysetFieldOpEq)
			if err != nil {
				return nil, err
			}
			and = append(and, expr)
		}

		f := fs[i]
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

		or = append(or, field.And(and...))
	}

	return field.Or(or...), nil
}
