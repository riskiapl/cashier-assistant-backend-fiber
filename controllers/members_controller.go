package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/riskiapl/fiber-app/services"
	"github.com/riskiapl/fiber-app/types"
)

type MemberController struct {
	service *services.MemberService
}

func NewMemberController(service *services.MemberService) *MemberController {
	return &MemberController{service: service}
}

func (c *MemberController) GetMembers(ctx *fiber.Ctx) error {
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	members, err := c.service.GetMembers(limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get members",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Members retrieved successfully",
		"data":    members,
	})
}

func (c *MemberController) GetMember(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid member ID",
			"error":   err.Error(),
		})
	}

	member, err := c.service.GetMember(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Member not found",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Member retrieved successfully",
		"data":    member,
	})
}

func (c *MemberController) UpdateMember(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid member ID",
			"error":   err.Error(),
		})
	}

	var req types.UpdateMemberRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	updatedMember, err := c.service.UpdateMember(uint(id), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update member",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Member updated successfully",
		"data":    updatedMember,
	})
}

func (c *MemberController) DeleteMember(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid member ID",
			"error":   err.Error(),
		})
	}

	if err := c.service.DeleteMember(uint(id)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete member",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Member deleted successfully",
	})
}
