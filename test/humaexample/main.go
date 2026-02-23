package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixlcrashr/go-pagetoken/test/humaexample/db/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultDSN = "host=127.0.0.1 port=5473 user=books password=books dbname=books sslmode=disable"

func main() {
	db, err := connectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get underlying sql.DB: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	srv := newServer(db)

	fmt.Println("Listening on 127.0.0.1:8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Listen("127.0.0.1:8080"); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	sig := <-quit
	fmt.Printf("\nReceived signal %s, shutting down gracefully...\n", sig)

	if err := srv.Shutdown(); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}

	fmt.Println("Server stopped.")
}

func connectToDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Migrator().DropTable(&model.Book{}); err != nil {
		return nil, fmt.Errorf("failed to drop tables: %v", err)
	}

	if err := db.AutoMigrate(&model.Book{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := seed(db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %v", err)
	}

	return db, nil
}
