package utils

//
//import (
//	"fmt"
//	"net/smtp"
//)
//
//func SendResetCode(email, code string) error {
//	from := "sizningemailingiz@gmail.com"
//	password := "abcd efgh ijkl mnop" // Gmail App password (bo‘sh joysiz yozing)
//
//	to := []string{email}
//	msg := []byte(fmt.Sprintf("Subject: Parolni tiklash\r\n\r\nSizning tiklash kodingiz: %s", code))
//
//	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
//	return smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
//}
