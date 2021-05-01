package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type outlookAuth struct {
	Username, Password string
}

func (auth *outlookAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(auth.Username), nil
}

func (auth *outlookAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(auth.Username), nil
		case "Password:":
			return []byte(auth.Password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}

func makeMessage(from string, to []string, subject string, body string) []byte {
	//TODO make everything https://stackoverflow.com/questions/50650719/golang-smtp-sending-empty-email
	From := fmt.Sprintf("From: <%s>\r\n", from)
	To := fmt.Sprintf("To: <%s>\r\n", strings.Join(to, ";"))
	Subject := fmt.Sprintf("Subject: %s\r\n", subject)
	Body := fmt.Sprintf("%s\r\n", body)
	msg := From + To + Subject + "\r\n" + Body
	return []byte(msg)
}

func loadUserData() (*outlookAuth, error) {
	file, err := ioutil.ReadFile("senderData.json")
	if err != nil {
		return nil, err
	}
	var auth outlookAuth
	err = json.Unmarshal(file, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func main() {
	senderAuth, err := loadUserData()
	if err != nil {
		fmt.Println("could not read password from password.txt")
	}

	to := []string{"example@gmail.com", "example@outlook.de"}

	smtpHost := "smtp-mail.outlook.com"
	smtpPort := "587"

	message := makeMessage(senderAuth.Username, to, "test-email", "Hello, this is an automated e-mail")

	err = smtp.SendMail(smtpHost+":"+smtpPort, senderAuth, senderAuth.Username, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}
