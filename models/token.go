package models

import "github.com/golang-jwt/jwt"

type (
	Claims struct {
		jwt.StandardClaims
		ID   uint   `gorm:"primaryKey"`
		Role string `json:"role"`
	}

	AccessToken struct {
		AccessToken string `json:"access"`
	}

	RefreshToken struct {
		RefreshToken string `json:"refresh"`
	}
)
