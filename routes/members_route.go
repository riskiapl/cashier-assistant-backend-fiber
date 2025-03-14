package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/riskiapl/fiber-app/controllers"
)

func SetupMemberRoutes(app *fiber.App, memberController *controllers.MemberController) {
	members := app.Group("/members")

	members.Get("/", memberController.GetMembers)
	members.Get("/:id", memberController.GetMember)
	members.Put("/:id", memberController.UpdateMember)
	members.Delete("/:id", memberController.DeleteMember)
}
