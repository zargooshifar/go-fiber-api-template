package auth

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers/auth"
)

func Routes(app *fiber.App) {
	app.Post("api/auth/checkuername", auth.CheckUserName)
	app.Post("api/auth/token", auth.Login)
	app.Post("api/auth/refresh", auth.Refresh)
	app.Post("api/auth/verifypin", auth.VerifyPin)
	app.Post("api/auth/completeregister", auth.CompleteRegister)
	app.Get("api/auth/logout", auth.Logout)
}
