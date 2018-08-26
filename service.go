package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bolsunovskyi/kubernetes-kong-certbot/kong"
)

type Service struct {
	certBotPath   string
	email         string
	kongClient KongClient
}

type KongClient interface {
	AddCertificate(cert, key, domain string) error
	GetCertificates() (*kong.CertListResponse, error)
}

func MakeService() *Service {
	return &Service{
		certBotPath:   "certbot",
		email:         "",
	}
}

func (s Service) StartRouter() error {
	r := gin.New()
	r.GET("/", func(context *gin.Context) {
		context.String(200, "ready")
	})
	r.Static("/.well-known/", "./static/.well-known/")
	r.POST("/certonly/:domain", s.certOnly)

	if err := r.Run(":8000"); err != nil {
		return err
	}

	return nil
}
