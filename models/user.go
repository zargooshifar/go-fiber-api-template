package models

import (
	"strings"
)

const (
	ADMIN = "admin"
	USER  = "user"
	CRM   = "crm"
)

type ROLES []string

func (list ROLES) Has(role string) bool {
	for _, target := range list {
		if strings.EqualFold(target, role) {
			return true
		}
	}
	return false
}

type (
	User struct {
		Base
		Username  string `json:"username"`
		Password  string `json:"-"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	}

	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Registration struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Password     string `json:"password"`
		Verification string `json:"verification_id"`
	}

	RegistrationError struct {
		Err          bool   `json:"error"`
		Verification string `json:"verification_id"`
		Password     string `json:"password"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}

	LoginUserName struct {
		Username string `json:"username"`
	}

	LoginPassword struct {
		Password string `json:"password"`
	}

	LoginErrors struct {
		Err           bool `json:"err"`
		UserExist     bool `json:"user_exist"`
		WrongPassword bool `json:"wrong_password"`
	}
)
