package core

import (
	"net/http"

	"github.com/gofiber/template/mustache/v2"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Engine struct {
	App *fiber.App
}

func NewEngine() Engine {
	app := fiber.New(fiber.Config{
		Views:       mustache.New("./views", ".mustache"),
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world from Fiber")
	})

	//rendering
	app.Get("/render", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World from Fiber!",
		})
	})

	app.Use(recover.New())

	app.Static("/", "./public")

	// api
	app.Get("/catalog", func(c *fiber.Ctx) error {
		if len(c.Query("param")) == 0 {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{
					"success": false,
					"error":   "catalog param not found",
				})
		}
		return c.JSON(&fiber.Map{
			"success": true,
			"data":    c.Query("param"),
		})
	})

	app.Use(cors.New())

	// middleware
	app.Use(func(c *fiber.Ctx) error {
		if c.Is("json") {
			return c.Next()
		}
		return c.SendString("Only JSON allowed!")
	})

	return Engine{
		App: app,
	}
}
