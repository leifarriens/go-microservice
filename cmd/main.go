package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gorm.io/driver/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/leifarriens/go-microservice/docs"
	"github.com/leifarriens/go-microservice/handler"
	"github.com/leifarriens/go-microservice/repository"
	"github.com/leifarriens/go-microservice/service"
	"github.com/leifarriens/go-microservice/utils"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("service.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	productRepository := repository.NewProductRepository(db)

	productService := service.NewProductService(&service.ProductServiceConfig{
		ProductRepository: productRepository,
	})

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(utils.Logger())

	e.Validator = utils.NewValidator()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handler.NewHandler(&handler.HandlerConfig{
		E:              e,
		ProductService: productService,
	})

	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	// Start server
	go func() {
		if err := e.Start(":" + "8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalln(err)
	}

	// shutdown db connection(s)
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
