package main

import (
	"context"
	"example/config"
	"example/entities"
	"example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	app := fiber.New()

	if _, err := config.DatabaseConnection(context.Background()); err != nil {
		log.Fatalf(" Gagal konek ke database: %v", err)
	}
	log.Println(" Database connected")

	if err := config.GormDB.AutoMigrate(&entities.Category{}); err != nil {
		log.Fatalf(" Gagal migrate tabel: %v", err)
	}
	log.Println(" Auto migration sukses")

	routes.CategoryRoutes(app)

	app.Listen(":3000")
}
