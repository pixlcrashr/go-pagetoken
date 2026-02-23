package gorm

import (
	"github.com/pixlcrashr/go-pagetoken"
	"gorm.io/gen"
)

func KeysetTokenCond(token *pagetoken.KeysetToken) gen.Condition {
	return gen.Cond()[0]
}
