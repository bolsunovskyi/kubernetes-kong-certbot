package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func (s Service) CertOnly(domain string) error {
	log.Println("certonly for", domain)

	cmd := exec.Command(s.certBotPath, "--config-dir", "./certbot/configs/", "--work-dir", "./certbot/work/",
		"--logs-dir", "./certbot/logs/", "certonly", "--webroot", "-w", "./static", "-d", domain, "-m", s.email,
		"--agree-tos", "--keep-until-expiring")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cert, err := os.Open("./certbot/configs/live/" + domain + "/fullchain.pem")
	if err != nil {
		return err
	}
	defer cert.Close()
	certBytes, _ := ioutil.ReadAll(cert)

	key, err := os.Open("./certbot/configs/live/" + domain + "/privkey.pem")
	if err != nil {
		return err
	}
	defer cert.Close()
	keyBytes, _ := ioutil.ReadAll(key)

	if err := s.kongClient.DeleteCertificate(domain); err != nil {
		log.Println(err)
	}

	if err := s.kongClient.AddCertificate(string(certBytes), string(keyBytes), domain); err != nil {
		return err
	}

	return nil
}
