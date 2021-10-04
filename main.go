package main

import (
	"flag"
	"log"
	"msgv2-back/database"
	"msgv2-back/handlers"
	"msgv2-back/routes/auth"
	"msgv2-back/routes/contact"
	"msgv2-back/routes/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {

	flag.Parse()

	database.ConnectDB()

	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	app.Use(recover.New())
	app.Use(logger.New())

	auth.Routes(app)
	users.Routes(app)
	contact.Routes(app)
	app.Use(handlers.NotFound)

	log.Fatal(app.Listen(*port))

	app.Listen(":3000")
}
