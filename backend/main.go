// main.go
package main

import (
	"log"

	"github.com/cs_student_uni/microservice_store/authService"
	"github.com/cs_student_uni/microservice_store/database"
	"github.com/cs_student_uni/microservice_store/storeService"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	authApp := fiber.New()

	authApp.Use(cors.New())
	authApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	storeApp := fiber.New()
	storeApp.Use(cors.New())
	storeApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	err := database.InitDB()
	if err != nil {
		log.Fatal("database error : ", err)
	}

	// Initialize default config

	api := authApp.Group("/api") // /api

	s1 := api.Group("/s1")
	s1.Post("/login", authService.Login)
	s1.Get("/users", authService.GetUsers)
	s1.Post("/signup", authService.Signup)

	apiStore := storeApp.Group("/api")

	s2 := apiStore.Group("/s2")                   // /api/s2
	s2.Get("/products", storeService.GetProducts) // /api/s2/list
	s2.Post("/createProduct", storeService.CreateProduct)
	s2.Post("/sellProduct", storeService.SellProduct)

	storeService.Init()
	authService.Init()

	// Start the HTTP servers for both services
	go func() {
		storeApp.Listen(":8080") // Store Service listens on port 8080
	}()
	authApp.Listen(":8081") // Authentication Service listens on port 8081

}
