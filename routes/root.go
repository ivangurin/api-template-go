package routes

import "github.com/gofiber/fiber/v2"

func GetRoot(c *fiber.Ctx) error {

	return c.SendString("Hello API!")

}