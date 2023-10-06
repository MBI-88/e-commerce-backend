package users

import (
	"fiber-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ConfigPath(router *fiber.App) *fiber.App {
	router.Post("/signup", signUp)
	router.Post("/signin", signIn)
	router.Post("/restorepassword", restorePassword)
	router.Post("/check-code", checkConfirmCode)
	router.Patch("/change-password", middlewares.JWTmiddleware, changePassword)
	router.Patch("/change-email", middlewares.JWTmiddleware, changeEmail)
	router.Get("/:pk<int>", middlewares.JWTmiddleware, getUser)
	router.Patch("/", middlewares.JWTmiddleware, updateUser)
	router.Delete("/", middlewares.JWTmiddleware, deleteUser)
	router.Patch("/change-role", middlewares.JWTmiddleware, changeRole)
	router.Get("/get-publishings/:user_id<int>", middlewares.JWTmiddleware, getPublishings)
	router.Post("/create-publishing", middlewares.JWTmiddleware, createPublishing)
	router.Patch("/update-publishing", middlewares.JWTmiddleware, updatePublishing)
	router.Delete("/delete-publishing", middlewares.JWTmiddleware, deletePublishing)
	router.Patch("/new-code", newConfiramtionCode)
	router.Get("/views", middlewares.JWTmiddleware, getViews)
	router.Post("/invalid-email", invalidEmail)
	router.Get("/comments", middlewares.JWTmiddleware, getComments)
	router.Get("/trolley/:pk<int>", middlewares.JWTmiddleware, getAllTrolley)
	router.Post("/trolley", middlewares.JWTmiddleware, createTrolley)
	router.Patch("/trolley", middlewares.JWTmiddleware, updateTrolley)
	router.Delete("/trolley", middlewares.JWTmiddleware, deleteTrolley)
	router.Delete("/all-trolley/:pk<int>", middlewares.JWTmiddleware, deleteAllTroley)
	return router
}
