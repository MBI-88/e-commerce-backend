package router

import (
	"fiber-backend/services/admin"
	"fiber-backend/services/publishing"
	"fiber-backend/services/users"

	"github.com/gofiber/fiber/v2"
)

// Router. Principal router for services
func Router(path *fiber.App) {
	path.Mount("/users", users.ConfigPath(fiber.New()))
	path.Mount("/publishings", publishing.ConfigPath(fiber.New()))
	path.Mount("/admin", admin.ConfigPath(fiber.New()))
}