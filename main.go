package main

import (
	"os"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var AXIOM *axiom.Client

type Todo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
}

func CreateAxiomClient() {
	AXIOM_TOKEN := os.Getenv("AXIOM_TOKEN")
	AXIOM_ORG_ID := os.Getenv("AXIOM_ORG_ID")

	client, err := axiom.NewClient(
		axiom.SetPersonalTokenConfig(AXIOM_TOKEN, AXIOM_ORG_ID),
	)
	if err != nil {
		panic("Could not create Axiom client")
	}

	AXIOM = client
}

func main() {
	app := fiber.New()

	ConnectDB()
	CreateAxiomClient()

	// Migrate the database
	DB.AutoMigrate(&Todo{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/todos", func(c *fiber.Ctx) error {
		var todos []Todo
		DB.Find(&todos)
		return c.JSON(todos)
	})

	app.Post("/todos", func(c *fiber.Ctx) error {
		var todo Todo
		if err := c.BodyParser(&todo); err != nil {
			return err
		}

		DB.Create(&todo)
		return c.JSON(todo)
	})

	app.Listen(":3000")
}
