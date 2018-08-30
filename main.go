package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bolsunovskyi/kubernetes-kong-certbot/kong"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalln("Unable to load env file")
		}
	}

}

type Config struct {
	CertBotPath     string `validate:"required"`
	Email           string `validate:"required"`
	AuthHeaderKey   string
	AuthHeaderValue string
	KongAdminPath   string `validate:"required"`
	Domain          string `validate:"required"`
	Listen          string `validate:"required"`
}

func main() {
	cnf := Config{Listen: ":8000"}
	if err := envconfig.Process("app", &cnf); err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(&cnf.Domain, "d", "", "domain for certificate")
	flag.Parse()

	if err := validator.New().Struct(cnf); err != nil {
		log.Fatalln(err)
	}

	var authHeader *kong.AuthHeader
	if cnf.AuthHeaderKey != "" && cnf.AuthHeaderValue != "" {
		authHeader = &kong.AuthHeader{Key: cnf.AuthHeaderKey, Value: cnf.AuthHeaderValue}
	}

	svc := MakeService(cnf.CertBotPath, cnf.Email, cnf.Listen,
		kong.Make(&http.Client{Timeout: time.Second * 30}, authHeader, cnf.KongAdminPath))

	if certs, err := svc.kongClient.GetCertificates(); err != nil {
		log.Fatalln(err)
	} else {
		log.Printf("certificates on server: %d\n", certs.Total)
	}

	go func() {
		if err := svc.StartRouter(); err != nil {
			log.Fatalln(err)
		}
	}()
	//
	if err := svc.CertOnly(cnf.Domain); err != nil {
		log.Fatalln(err)
	}
}
