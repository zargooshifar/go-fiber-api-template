package sms

import (
	"log"
	"math/rand"
	"msgv2-back/database"
	"msgv2-back/models"
	"time"
)


func SendPin(number string) (string, error) {
	pin := 10000 + rand.Intn(89999)

	//delete previous verifications
	database.DB.Where(&models.VerificationSMS{Number: number}).Delete(&models.VerificationSMS{})

	verification := new(models.VerificationSMS)

	verification.Pin = pin
	verification.Number = number
	verification.Expire = time.Now().Add(2 * time.Minute).Unix()

	if err := database.DB.Create(&verification).Error; err != nil {
		return "", err
	}

	//TODO: implement sending sms
	//send sms
	log.Println("sending pin: %s", pin)
	return verification.ID.String(), nil
}
