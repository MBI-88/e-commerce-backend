package publishing

import (
	"fiber-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ConfigPath(route *fiber.App) *fiber.App {
	route.Get("/", getAll)
	route.Get("/top4-sellers", getTopSeller)
	route.Get("/top9", getTop9Publishings)
	route.Get("/filters", getPublishingsByFilters)
	route.Get("/:pk<int>", getPublishing)
	route.Get("/ccount", getCountByCategories)
	route.Get("/top9-category/:category<string>", getTopPublishingsByCategory)
	route.Get("/top9-subcategory/:subcategory<string>", getTopPublishingsBySubCategory)
	route.Get("/scount/:category<string>", getCountBySubCategories)
	route.Get("/:category<string>", getPublishingsByCategories)
	route.Get("/:category<string>/:subcategory<string>", getPublishingsBySubCategories)
	route.Post("/views", createViews)
	route.Post("/comment", middlewares.JWTmiddleware, createComments)
	return route
}
