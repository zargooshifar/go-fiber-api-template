package users

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/admin/users", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.GetItems(models.User{}))
	app.Get("api/admin/user", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.GetItem(models.User{}))
	app.Put("api/admin/user", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.CreateItem(models.User{}))
	app.Post("api/admin/user", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.UpdateItem(models.User{}))
	app.Delete("api/admin/user", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.DeleteItem(models.User{}))
}
