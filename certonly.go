package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"log"
	"os"
	"io/ioutil"
)

func (s Service) certOnly(context *gin.Context) {
	domain := context.Param("domain")
	if domain == "" {
		context.Status(400)
		return
	}

	log.Println("certonly for", domain)

	cmd := exec.Command(s.certBotPath, "--config-dir", "./certbot/configs/", "--work-dir", "./certbot/work/",
		"--logs-dir", "./certbot/logs/", "certonly", "--webroot", "-w", "./static", "-d", domain, "-m", s.email,
		"--agree-tos", "--keep-until-expiring")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		context.String(500, err.Error())
		return
	}

	context.String(200, "certonly for %s\n", domain)

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Println(err)
			return
		}

		cert, err := os.Open("./certbot/configs/live/"+domain+"/fullchain.pem")
		if err != nil {
			log.Println(err)
			return
		}
		defer cert.Close()
		certBytes, _ := ioutil.ReadAll(cert)

		key, err := os.Open("./certbot/configs/live/"+domain+"/privkey.pem")
		if err != nil {
			log.Println(err)
			return
		}
		defer cert.Close()
		keyBytes, _ := ioutil.ReadAll(key)
		//TODO: delete by SNIS before add
		if err := s.kongClient.AddCertificate(string(certBytes), string(keyBytes), domain); err != nil {
			log.Println(err)
			return
		}
	}()
}
