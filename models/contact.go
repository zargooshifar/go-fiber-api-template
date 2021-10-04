package models



type Contact struct {
	Base
	Name			string `json:"name"`
	Phone			string `json:"phone"`
}
