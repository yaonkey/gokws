package core

import (
	"net/http"
	"time"

	"github.com/gofiber/template/mustache/v2"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/golang-jwt/jwt/v5"
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

	cfg := swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "KWS API",
	}

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(swagger.New(cfg))

	app.Post("/login", loginHandler)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
	app.Get("/restricted", restrictedHandler)

	// static files
	app.Static("/js", "./public/js")
	app.Static("/css", "./public/css")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world from Fiber")
	})

	//rendering
	app.Get("/render", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World from Fiber!",
		}, "layouts/main")
	})

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

func loginHandler(c *fiber.Ctx) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	// TODO: Throws Unauthorized error
	if login != "admin" || password != "admin" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":       "Admin",
		"is_admin":   true,
		"expired_at": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("env_secret")) // todo
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

// todo change it
func restrictedHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
