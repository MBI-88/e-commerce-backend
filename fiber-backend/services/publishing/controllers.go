package publishing

import (
	"fiber-backend/helpers"
	"fiber-backend/models"

	"github.com/gofiber/fiber/v2"
)

// views

// Return top 9 publishings
func getTop9Publishings(ctx *fiber.Ctx) error {
	var views models.Views
	pubs, err := views.GetViews()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"top9": pubs})
}

// Return a publishing with its information
func getPublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	pk, err := ctx.ParamsInt("pk", 1)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := pub.GetPublishing(pk); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishing": pub})
}

// Get publishings by categories
func getPublishingsByCategories(ctx *fiber.Ctx) error {
	var pub models.Publishings
	params := ctx.Params("category", "")
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	pubs, err := pub.GetPublishingsByCategories(params, offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings": pubs})
}

// Get publishings by subcategory having category
func getPublishingsBySubCategories(ctx *fiber.Ctx) error {
	var pub models.Publishings
	cate := ctx.Params("category", "")
	sub := ctx.Params("subcategory", "")
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	pubs, err := pub.GetPublishingsBySubCategory(offset, size, cate, sub)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings": pubs})
}

// Get publishings by filters having category and subcategory
func getPublishingsByFilters(ctx *fiber.Ctx) error {
	var pub models.Publishings
	var queryStruct helpers.ParseQuery
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	if err := ctx.QueryParser(&queryStruct); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	pubs, err := pub.GetPublishingsByFilters(queryStruct, offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings": pubs})
}

// Get publishings count by category
func getCountByCategories(ctx *fiber.Ctx) error {
	var pub models.Publishings
	result, err := pub.GetCountByCategories()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"counts": result})
}

// Get publishings count by sub categories having category
func getCountBySubCategories(ctx *fiber.Ctx) error {
	var pub models.Publishings
	params := ctx.Params("category", "Technology")
	result, err := pub.GetCountBySubCategories(params)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"counts": result})
}

// Create views. Save publishings views
func createViews(ctx *fiber.Ctx) error {
	var view models.Views
	if err := ctx.BodyParser(&view); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := view.SaveView(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Create view error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Create comments. Make comments to a publishing
func createComments(ctx *fiber.Ctx) error {
	var com models.Comments
	if err := ctx.BodyParser(&com); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(com); errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(errors)
	}
	if err := com.CreateComent(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Create comment error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful"})
}

// Get publishings by category with top views
func getTopPublishingsByCategory(ctx *fiber.Ctx) error {
	var v models.Views
	cate := ctx.Params("category", "")
	from := ctx.Query("start_date", "")
	to := ctx.Query("end_date", "")
	top9, err := v.GetTopByCategory(cate,from, to)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal error"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"top9": top9})
}

// Get publishings by sub-category with top views
func getTopPublishingsBySubCategory(ctx *fiber.Ctx) error {
	var v models.Views
	sub := ctx.Params("subcategory", "")
	from := ctx.Query("start_date", "")
	to := ctx.Query("end_date", "")
	top9, err := v.GetTopBySubCategory(sub,from,to)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal error"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"top9": top9})
}

// Get all publishings
func getAll(ctx *fiber.Ctx) error {
	var pub models.Publishings
	var queryStruct helpers.ParseQuery
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	if err := ctx.QueryParser(&queryStruct); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	pubs, err := pub.GetAll(queryStruct, offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings": pubs})
}

// Get top 4 sellers
func getTopSeller(ctx *fiber.Ctx) error {
	var c models.Comments
	sellers, err := c.GetTopSellerByRating()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"sellers": sellers})
}
