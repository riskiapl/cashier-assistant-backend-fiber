package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/controllers"
)

func MemberRoutes(app *fiber.App) {
	members := app.Group("/members")

	memberController := controllers.NewMemberController()

	members.Get("/", memberController.GetMembers)
	members.Get("/:id", memberController.GetMember)
	members.Put("/:id", memberController.UpdateMember)
	members.Delete("/:id", memberController.DeleteMember)
}
