package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/sms"
	"msgv2-back/models"
	"time"
)

func Login(c *fiber.Ctx) error {
	l := new(models.Login)

	if err := c.BodyParser(l); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	user := new(models.User)

	if count := database.DB.Where(&models.User{Username: l.Username}).First(&user).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.USER_NOT_EXIST,
		})
	}

	//check password
	if !utils.VerifyPassword(user.Password, l.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.ERROR_WRONG_PASSWORD,
		})
	}

	accessToken, refreshToken := utils.GenerateTokens(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  accessToken,
		"refresh": refreshToken,
	})

}

func CheckUserName(c *fiber.Ctx) error {
	username := new(models.LoginUserName)

	if err := c.BodyParser(username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}
	exists := (database.DB.Where(&models.User{Username: username.Username}).First(&models.User{}).RowsAffected > 0)
	if !exists {
		verification_id, err := sms.SendPin(username.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"verification_id": verification_id,
			"exists":          exists,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exists": exists,
	})
}

func VerifyPin(c *fiber.Ctx) error {
	max_attempts := 5
	pin := new(models.Pin)
	if err := c.BodyParser(pin); err != nil {
		return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	verification := new(models.VerificationSMS)

	exists := (database.DB.Where(&models.VerificationSMS{ID: pin.ID}).First(&verification).RowsAffected > 0)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_VERIFICATION_ID,
		})
	}

	if verification.Expire < time.Now().Unix() || verification.Attempts > max_attempts {
		database.DB.Delete(verification)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.EXPIRE_PIN,
		})
	}

	if verification.Pin == pin.Pin {
		verification.Confirm = true
		database.DB.Save(verification)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":         true,
			"message":         "",
			"verification_id": verification.ID.String(),
		})
	} else {
		verification.Attempts += 1
		database.DB.Save(verification)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_PIN,
		})
	}

}

func CompleteRegister(c *fiber.Ctx) error {

	reg := new(models.Registration)
	if err := c.BodyParser(reg); err != nil {
		return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	log.Println(reg.Verification)
	log.Println(reg.FirstName)
	log.Println(reg.LastName)
	log.Println(reg.Password)
	// validate if the email, username and password are in correct format
	regErrors := utils.ValidateRegister(reg)
	if regErrors.Err {
		return c.Status(fiber.StatusBadRequest).JSON(regErrors)
	}

	verification := new(models.VerificationSMS)
	if count := database.DB.Where(&models.VerificationSMS{ID: uuid.MustParse(reg.Verification)}).First(&verification).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_EXIST,
		})
	}

	if !verification.Confirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_CONFIRMED,
		})
	}

	if count := database.DB.Where(&models.User{Username: verification.Number}).First(new(models.User)).RowsAffected; count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.USER_EXIST,
		})
	}

	user := new(models.User)

	// Hashing the password with a random salt
	password := []byte(reg.Password)
	hashedPassword, err := utils.GenerateHashPassword(password)

	if err != nil {
		panic(err)
	}
	user.Username = verification.Number
	user.FirstName = reg.FirstName
	user.LastName = reg.LastName
	user.Password = string(hashedPassword)
	user.Role = "user"
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}

	//delete verification...
	database.DB.Delete(verification)

	log.Println("user created!")

	// setting up the authorization cookies
	accessToken, refreshToken := utils.GenerateTokens(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  accessToken,
		"refresh": refreshToken,
	})
}

func Refresh(c *fiber.Ctx) error {
	r := new(models.RefreshToken)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}
	access, refresh, error := utils.RefreshTokens(r)
	if len(error) > 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  access,
		"refresh": refresh,
	})
}

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user")
	id := c.Locals("id")
	log.Println("logout")
	log.Println(user)
	log.Println(id)
	return c.SendStatus(fiber.StatusOK)
}
