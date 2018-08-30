package main

import (
	"github.com/bolsunovskyi/kubernetes-kong-certbot/kong"
	"github.com/gin-gonic/gin"
)

type Service struct {
	certBotPath  string
	email        string
	routerListen string
	kongClient   KongClient
}

type KongClient interface {
	AddCertificate(cert, key, domain string) error
	GetCertificates() (*kong.CertListResponse, error)
	UpdateOrCreateCertificate(host, cert, key string) error
	DeleteCertificate(host string) error
}

func MakeService(certBotPath, email, routerListen string, kongClient KongClient) *Service {
	return &Service{
		certBotPath:  certBotPath,
		email:        email,
		kongClient:   kongClient,
		routerListen: routerListen,
	}
}

func (s Service) StartRouter() error {
	r := gin.New()
	r.GET("/", func(context *gin.Context) {
		context.String(200, "ready")
	})
	r.Static("/.well-known/", "./static/.well-known/")

	if err := r.Run(s.routerListen); err != nil {
		return err
	}

	return nil
}
