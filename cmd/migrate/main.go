package main

import (
	"log"

	"github.com/cryptoPickle/go-ddd-example/config"
	postgres "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress"
	postgres_dto "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/dto"
	"github.com/cryptoPickle/go-ddd-example/migrations/seed"
)

func main() {
	c := config.NewConfig()

	db, err := postgres.NewPostgressConnection(c.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}
	db.Debug()

	db.AutoMigrate(
		&postgres_dto.Item{},
		&postgres_dto.Payout{},
		&postgres_dto.Sale{},
		&postgres_dto.BatchPayout{},
	)

	itemSeeds := seed.ItemSeeds(1000)
	itemSeeds.RunAll(db)
}
