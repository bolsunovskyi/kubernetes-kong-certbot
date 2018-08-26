package main

import (
	"log"
	"os"
	"github.com/bolsunovskyi/kubernetes-kong-certbot/kong"
	"net/http"
)

func main() {
	svc := MakeService()

	if certbotPath := os.Getenv("CERTBOT_PATH"); certbotPath != "" {
		svc.certBotPath = certbotPath
	}

	if email := os.Getenv("EMAIL"); email != "" {
		svc.email = email
	}

	authHeaderKey := os.Getenv("AUTH_HEADER_KEY")
	authHeaderValue := os.Getenv("AUTH_HEADER_VALUE")
	kongAdminPath := os.Getenv("KONG_ADMIN_PATH")
	if kongAdminPath == "" {
		log.Fatalln("kong path is not specified")
	}


	if authHeaderKey != "" && authHeaderValue != "" {
		svc.kongClient = kong.Make(http.DefaultClient, &kong.AuthHeader{Value: authHeaderValue, Key: authHeaderKey},
		kongAdminPath)
	} else {
		svc.kongClient = kong.Make(http.DefaultClient, nil, kongAdminPath)
	}

	if certs, err := svc.kongClient.GetCertificates(); err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("certificates on server: %d\n", certs.Total)
	}

	if err := svc.StartRouter(); err != nil {
		log.Fatalln(err)
	}
}
