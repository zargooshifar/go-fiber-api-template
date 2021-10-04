package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"os"
	"strings"
	"time"
)

var access_key = []byte(os.Getenv("ACCESS_KEY"))
var refresh_key = []byte(os.Getenv("REFRESH_KEY"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(user *models.User) (string, string) {
	claim, accessToken := GenerateAccessClaims(user)
	refreshToken := GenerateRefreshClaims(claim)

	return accessToken, refreshToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(user *models.User) (*models.Claims, string) {

	t := time.Now()
	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: t.Add(15 * time.Minute).Unix(),
			Subject:   "access",
			IssuedAt:  t.Unix(),
		},
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(access_key)
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(cl *models.Claims) string {
	result := database.DB.Where(&models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: cl.Issuer,
		},
	}).Find(&models.Claims{})

	// checking the number of refresh tokens stored.
	// If the number is higher than 3, remove all the refresh tokens and leave only new one.
	if result.RowsAffected > 3 {
		database.DB.Where(&models.Claims{
			StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer},
		}).Delete(&models.Claims{})
	}

	t := time.Now()
	refreshClaim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: t.Add(30 * 24 * time.Hour).Unix(),
			Subject:   "refresh",
			IssuedAt:  t.Unix(),
		},
		Role: cl.Role,
	}

	// create a claim on DB
	database.DB.Create(&refreshClaim)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(refresh_key)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}

// SecureAuth returns a middleware which secures all the private routes
func Secure(roles models.ROLES) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		reqToken := c.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		accessToken := splitToken[1]

		claims := new(models.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return access_key, nil
			})

		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"general": "Token Expired",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// not a token!

				return c.SendStatus(fiber.StatusForbidden)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return c.SendStatus(fiber.StatusUnauthorized)
			} else {
				// cannot handle this token

				return c.SendStatus(fiber.StatusForbidden)
			}
		}

		user := new(models.User)
		if (database.DB.Where(models.User{
			Base: models.Base{
				ID: uuid.MustParse(claims.Issuer),
			},
		}).Find(&user).Error != nil) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.USER_NOT_EXIST,
			})
		}

		if roles.Has(claims.Role) {
			c.Locals("user", user)
			c.Locals("id", user.ID.String())
			return c.Next()
		}

		return c.SendStatus(fiber.StatusForbidden)

	}
}

func RefreshTokens(r *models.RefreshToken) (string, string, string) {
	claims := new(models.Claims)

	token, err := jwt.ParseWithClaims(r.RefreshToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return refresh_key, nil
		})

	if err != nil {
		return "", "", errors.REFRESH_TOKEN_UNVALID
	}

	if token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return "", "", errors.REFRESH_TOKEN_EXPIRED

		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			// not a token!
			return "", "", errors.REFRESH_TOKEN_NOT_A_TOKEN
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return "", "", errors.REFRESH_TOKEN_UNVALID

		} else {
			return "", "", errors.CANT_HANDLE_REFRESH_TOKEN
		}
	}

	//check if refresh is in db
	dbtoken := new(models.Claims)
	count := database.DB.Where(models.Claims{ID: claims.ID}).Find(&dbtoken).RowsAffected
	if count == 0 {
		return "", "", errors.REFRESH_TOKEN_NOT_EXIST
	}

	//delete previous refresh key
	database.DB.Delete(dbtoken)
	user := new(models.User)
	database.DB.Where(models.User{Base: models.Base{
		ID: uuid.MustParse(claims.Issuer),
	}}).Find(&user)
	accessToken, refreshToken := GenerateTokens(user)

	return accessToken, refreshToken, ""

}

func GenerateHashPassword(password []byte) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword(
		password,
		bcrypt.MinCost+5,
	)

	return hash, err
}

func VerifyPassword(hashPassword string, password string) bool {

	byteHash := []byte(hashPassword)
	bytePlain := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)

	if err != nil {
		return false
	}

	return true
}
