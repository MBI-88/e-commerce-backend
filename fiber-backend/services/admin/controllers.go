package admin

import (
	"fiber-backend/helpers"
	"fiber-backend/middlewares"
	"fiber-backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Login admin
func signIn(ctx *fiber.Ctx) error {
	var user models.Users
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Error data"})
	}
	if errors := helpers.ValidateStruct(user); errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(errors)
	}
	if err := user.LogUser(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	token, err := middlewares.JWTgenerate(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error Token"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"session_token": fmt.Sprintf("Bearer %s", token),
		"email":         user.Email,
		"username":      user.Username,
		"image":         user.Image,
	})
}

//****************************************************
//**************Users section*******************
//****************************************************

// Get users array
func getUsers(ctx *fiber.Ctx) error {
	var user models.Users
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	users, err := user.GetUsers(offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found!"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

// Update an user
func updateUser(ctx *fiber.Ctx) error {
	var user models.Users
	if err := user.EditUsers(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Update user error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "User updated successful!"})
}

// Create an user
func createUser(ctx *fiber.Ctx) error {
	var user models.Users
	if err := user.MakeUser(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Created user error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "User created successful!"})
}

// Delete an user
func deleteUser(ctx *fiber.Ctx) error {
	var user models.Users
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := user.DeleteUser(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Delete user error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "User deleted successful!"})
}

// Search user by filters
func searchUser(ctx *fiber.Ctx) error {
	var user models.Users
	username := ctx.Query("username", "")
	email := ctx.Query("email", "")
	role := ctx.Query("role", "")
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	users, err := user.SearchUsers(username, email, role, offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Search user error"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

// Return users by role
func usersByRole(ctx *fiber.Ctx) error {
	var user models.Users
	result, err := user.GetUsersByRole()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"count": result})
}

//****************************************************
//**************End section***************************
//****************************************************

//****************************************************
//**************Publishings section*******************
//****************************************************

// Search publishing using a key or return all publishings
func searchPublishings(ctx *fiber.Ctx) error {
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found!"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings":pubs})
}

// Create publishing 
func createPublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings 
	if err := pub.CreatePublishing(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "publishin error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Publishing created succesful!"})
}

// Update publishing
func updatePublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	if err := pub.UpdatePublishing(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "publishin error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Publishing updated succesful!"})
}

// Delete publishing
func deletePublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	if err := ctx.BodyParser(&pub); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := pub.DeletePublishing(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "publishin error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Publishing deleted succesful!"})
}

// Get total publishings by category
func publishingsCountByCategory(ctx *fiber.Ctx) error {
	var pub models.Publishings
	result, err := pub.GetCountByCategories()
	if err != nil {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"counts": result})
}

// Get total publishings by subcategory
func publishingsCountBySubcategory(ctx *fiber.Ctx) error {
	var pub models.Publishings 
	params := ctx.Params("category", "Technology")
	result, err := pub.GetCountBySubCategories(params)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"counts": result})
}

//****************************************************
//**************End section***************************
//****************************************************


//****************************************************
//**************Views section*************************
//****************************************************

// Get views
func getViews(ctx *fiber.Ctx) error {
	var v models.Views
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	cat := ctx.Query("category", "")
	sub := ctx.Query("subcategory", "")
	pubId := ctx.QueryInt("publishing_id", 0)
	offset, size := helpers.Page(page, pageSize)
	vs, err := v.GetViewsData(offset,size,pubId,cat,sub)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found!"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"views":vs})
}

// Create view
func saveView(ctx *fiber.Ctx) error {
	var v models.Views
	if err := ctx.BodyParser(&v); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(v); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if err := v.SaveView(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message":"Create view error!"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Action successful!"})
}


// Delete view
func deleteView(ctx *fiber.Ctx) error {
	var v models.Views
	if err := ctx.BodyParser(&v); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(v); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if err := v.DeleteView(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message":"Delete view error!"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message":"Action successful!"})
}

//****************************************************
//**************End section***************************
//****************************************************