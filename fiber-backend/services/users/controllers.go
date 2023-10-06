package users

import (
	"fiber-backend/helpers"
	"fiber-backend/middlewares"
	"fiber-backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Views

// Return an user with all information about him
func getUser(ctx *fiber.Ctx) error {
	var user models.Users
	param, err := ctx.ParamsInt("pk")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := user.GetUser(param); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(user)
}

// Log an user into the system
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
		"user_id":       user.ID,
		"image":         user.Image,
		"username":      user.Username,
		"session_token": fmt.Sprintf("Bearer %s", token)})
}

// Set up an user into the system
func signUp(ctx *fiber.Ctx) error {
	var user models.Users
	var c models.ConfirmationCode
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(user); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if err := user.CreateUser(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Create user error"})
	}
	if err := c.GenerateCode(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"code": c.Code,
		"code_id": c.ID, "user_id": user.ID, "email": user.Email})
}

// Update partial user's data
func updateUser(ctx *fiber.Ctx) error {
	var user models.Users
	if err := user.UpdateUser(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Update user error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(user)
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

// Change user's password
func changePassword(ctx *fiber.Ctx) error {
	var user models.Users
	var c models.ConfirmationCode
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(user); errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(errors)
	}
	if err := user.ChangePassword(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Change password error"})
	}
	if err := c.GenerateCode(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"code": c.Code,
		"code_id": c.ID, "user_id": user.ID, "email": user.Email})
}

// Change user's email
func changeEmail(ctx *fiber.Ctx) error {
	var user models.Users
	var c models.ConfirmationCode
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(user); errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(errors)
	}
	if err := user.ChangeEmail(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Change email error"})
	}
	if err := c.GenerateCode(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"new_email": user.Email, "user_id": user.ID,
		"code": c.Code, "code_id": c.ID})
}

// Check a confirmation code
func checkConfirmCode(ctx *fiber.Ctx) error {
	var c models.ConfirmationCode
	var user models.Users
	if err := ctx.BodyParser(&c); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := c.CheckCode(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	user = c.Users
	if err := user.ActiveUser(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Activate user error"})
	}
	token, err := middlewares.JWTgenerate(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Token error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"user_id": user.ID, "session_token": fmt.Sprintf("Bearer %s", token)})
}

// Restore password
func restorePassword(ctx *fiber.Ctx) error {
	var user models.Users
	var c models.ConfirmationCode
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if errors := helpers.ValidateStruct(user); errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(errors)
	}
	if err := user.GetUserbyEmail(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Get user error"})
	}
	if err := c.GenerateCode(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"code_id": c.ID,
		"code": c.Code, "email": user.Email, "user_id": user.ID})
}

// Change role
func changeRole(ctx *fiber.Ctx) error {
	var user models.Users
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := user.ChangeRole(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Role error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Re-send confirmation code
func newConfiramtionCode(ctx *fiber.Ctx) error {
	var c models.ConfirmationCode
	if err := ctx.BodyParser(&c); err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := c.UpdateCode(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Confirmation code error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"code_id": c.ID,
		"code": c.Code, "user_id": c.Users.ID, "email": c.Users.Email})
}

// Invalid user. Delete user if email is invalid
func invalidEmail(ctx *fiber.Ctx) error {
	var c models.ConfirmationCode
	var u models.Users
	if err := ctx.BodyParser(&c); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := c.DeleteCode(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Deleted code error"})
	}
	if err := u.DeleteInvalidUser(c.UserId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Deleted user error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "User deleted successful"})
}

//****************************************************
//**************Publishings section*******************
//****************************************************

// Create an item only used by sellers
func createPublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	if err := pub.CreatePublishing(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Create publishing error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Update an item only used by sellers
func updatePublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	if err := pub.UpdatePublishing(ctx); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Update publishing error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})

}

// Delete an item only used by sellers
func deletePublishing(ctx *fiber.Ctx) error {
	var pub models.Publishings
	if err := ctx.BodyParser(&pub); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	if err := pub.DeletePublishing(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Delete publishing error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})

}

// Get publishings by user
func getPublishings(ctx *fiber.Ctx) error {
	var p models.Publishings
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	pk, err := ctx.ParamsInt("user_id", -1)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	pubs, err := p.GetPublishingBySeller(offset, size, pk)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Get publishing error"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"publishings": pubs})
}

// Get views. Return views by publishing
func getViews(ctx *fiber.Ctx) error {
	var v models.Views
	startDate := ctx.Query("start_date", "")
	endDate := ctx.Query("end_date", "")
	pk := ctx.QueryInt("pk", 1)
	resp, err := v.GetViewsByPublishing(startDate, endDate, pk)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Get views error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"views": resp})
}

// Get comments by publishing
func getComments(ctx *fiber.Ctx) error {
	var com models.Comments
	page := ctx.QueryInt("page", 0)
	pageSize := ctx.QueryInt("page_size", 0)
	offset, size := helpers.Page(page, pageSize)
	pubID := ctx.QueryInt("pub", 0)
	comments, err := com.GetComments(pubID, offset, size)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"comments": comments})
}

//****************************************************
//**************Trolley section*******************
//****************************************************

// Create item in Trolley
func createTrolley(ctx *fiber.Ctx) error {
	var t models.Trolley
	if err := ctx.BodyParser(&t); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}
	if err := t.AddToTrolley(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Add item error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Update item in Trolley
func updateTrolley(ctx *fiber.Ctx) error {
	var t models.Trolley
	if err := ctx.BodyParser(&t); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}
	if err := t.UpdateTrolley(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Update item error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Delete item in Trolley
func deleteTrolley(ctx *fiber.Ctx) error {
	var t models.Trolley
	if err := ctx.BodyParser(&t); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}
	if err := t.DeleteFromTrolley(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Delete item error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Delete all Trolley
func deleteAllTroley(ctx *fiber.Ctx) error {
	var t models.Trolley
	pk, err := ctx.ParamsInt("pk", 0)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}
	if err := t.DeleteAll(pk); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Delete all error"})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Successful!"})
}

// Get all Trolley
func getAllTrolley(ctx *fiber.Ctx) error {
	var t models.Trolley
	pk, err := ctx.ParamsInt("pk", 0)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}
	trolley, err := t.GetAllTrolley(pk)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Get trolley error"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"items": trolley})
}
