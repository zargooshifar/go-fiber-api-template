package contact

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/admin/contacts", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.GetItems(models.Contact{}))
	app.Get("api/admin/contact", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.GetItem(models.Contact{}))
	app.Put("api/contact", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.CreateItem(models.Contact{}))
	app.Post("api/admin/contact", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.UpdateItem(models.Contact{}))
	app.Delete("api/admin/contact", utils.Secure(models.ROLES{models.ADMIN, models.CRM}), handlers.DeleteItem(models.Contact{}))
}

