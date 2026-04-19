package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"nisst/internal/config"
	"nisst/internal/database"
	"nisst/internal/handler"
	"nisst/internal/middleware"
	"nisst/internal/repository"
	"nisst/internal/service"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)
	repos := repository.NewRegistry(db)
	svcs := service.NewRegistry(repos)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.CORSOrigins}))
	app.Use(middleware.RequestLogger())

	handler.Register(app, svcs)

	log.Printf("server listening on %s:%s", cfg.Host, cfg.Port)
	log.Fatal(app.Listen(cfg.Host + ":" + cfg.Port))
}
