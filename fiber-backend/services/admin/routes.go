package admin

import (
	"fiber-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ConfigPath(routes *fiber.App) *fiber.App {
	routes.Post("/signin", signIn)
	routes.Get("/get-users", middlewares.JWTadminMiddleware, getUsers)
	routes.Post("/create-users", middlewares.JWTadminMiddleware, createUser)
	routes.Patch("/update-users", middlewares.JWTadminMiddleware, updateUser)
	routes.Delete("/delete-users", middlewares.JWTadminMiddleware, deleteUser)
	routes.Get("/search-users", middlewares.JWTadminMiddleware, searchUser)
	routes.Get("/get-publishings", middlewares.JWTadminMiddleware, searchPublishings)
	routes.Post("/create-publishing", middlewares.JWTadminMiddleware, createPublishing)
	routes.Patch("/update-publishing", middlewares.JWTadminMiddleware, updatePublishing)
	routes.Delete("/delete-publishing", middlewares.JWTadminMiddleware, deletePublishing)
	routes.Get("/get-byrole", middlewares.JWTadminMiddleware, usersByRole)
	routes.Get("/get-bycategory", middlewares.JWTadminMiddleware, publishingsCountByCategory)
	routes.Get("/get-bysubcategory/:category", middlewares.JWTadminMiddleware, publishingsCountBySubcategory)
	routes.Get("/get-views", middlewares.JWTadminMiddleware, getViews)
	routes.Post("/save-views", middlewares.JWTadminMiddleware, saveView)
	routes.Delete("/delete-views", middlewares.JWTadminMiddleware, deleteView)
	return routes
}
