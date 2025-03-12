package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/cryptoPickle/go-ddd-example/config"
	"github.com/cryptoPickle/go-ddd-example/internal/application/services"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/shared_services"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/usecase"
	infrastructure "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/apis"
	postgres "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/repositories"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger/zap"
	restapi "github.com/cryptoPickle/go-ddd-example/internal/interface/restapi"
)

func main() {
	c := config.NewConfig()
	logging, err := zap.NewZapLogger(&zap.ZapConfig{
		LogLevel: logger.Info,
		Env:      "prod",
		LogFile:  "logfile",
	})
	if err != nil {
		log.Fatalf("failed to create logger instance %v", err)
	}

	db, err := postgres.NewPostgressConnection(c.PostgresDSN())
	if err != nil {
		logging.Fatalf("failed to connect to db %v", err)
	}

	itemRepo := repositories.NewPostgresItemRepository(db)
	payoutRepo := repositories.NewPostgresPayoutRepository(db)

	exchageRateProvider := new(infrastructure.ExchangeRateProvider)
	exchangeRateService := shared_services.NewCurrencyConverter(exchageRateProvider)

	payoutCalculatorUsecase := usecase.NewPayoutCalculator(exchangeRateService)

	withLogger := services.NewPayoutServiceWithLogger
	payoutService := withLogger(services.NewPayoutService(itemRepo, payoutRepo, payoutCalculatorUsecase), logging)

	g := gin.Default()

	g.StaticFile("/docs", "./docs/swagger.json")
	g.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:3000/docs"),
	))

	restapi.NewPayoutController(g, payoutService, logging)

	g.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err := g.Run(":3000"); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
