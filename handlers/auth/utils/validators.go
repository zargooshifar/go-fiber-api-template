package utils

import (
	valid "github.com/asaskevich/govalidator"
	"msgv2-back/errors"
	"msgv2-back/models"
	"regexp"
)

func IsEmpty(str string) bool {
	if valid.HasWhitespaceOnly(str) && str != "" {
		return true
	}
	return false
}

// ValidateRegister func validates the body of user for registration
func ValidateRegister(u *models.Registration) *models.RegistrationError {
	e := &models.RegistrationError{}
	e.Err, e.Verification = IsEmpty(u.Verification), errors.EMPTY_VERIFICATION
	e.Err, e.FirstName = IsEmpty(u.FirstName), errors.EMPTY_FIRST_NAME
	e.Err, e.LastName = IsEmpty(u.LastName), errors.EMPTY_LAST_NAME

	re := regexp.MustCompile("\\d") // regex check for at least one integer in string
	if !(len(u.Password) >= 8 && valid.HasLowerCase(u.Password) && valid.HasUpperCase(u.Password) && re.MatchString(u.Password)) {
		e.Err, e.Password = true, "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"
	}
	return e
}
