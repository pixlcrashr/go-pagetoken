package main

import (
	"flag"

	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// gen-dao generates type-safe GORM query builders from the PostgreSQL database.
//
// Usage:
//
//	go run ./gen-dao/ [-dsn <dsn>] [-o <output-dir>]
//
// The compose stack must be running and the server started at least once
// (so that the schema exists) before invoking this tool.
func main() {
	dsn := flag.String(
		"dsn",
		"host=127.0.0.1 port=5473 user=books password=books dbname=books sslmode=disable",
		"PostgreSQL DSN",
	)
	outPath := flag.String("o", "../../db/model/dao", "Output directory for generated DAO code")
	flag.Parse()

	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       *outPath,
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	g.ApplyBasic(
		model.Book{})

	g.Execute()
}
