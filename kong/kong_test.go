package kong

import "testing"

func TestClient_AddCertificate(t *testing.T) {
	c := Make(nil, &AuthHeader{
		Key:   "###",
		Value: "###",
	}, "###")
	if err := c.AddCertificate("foo", "bar", "###"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCertificates(t *testing.T) {
	c := Make(nil, &AuthHeader{
		Key:   "###",
		Value: "###",
	}, "###")
	certs, err := c.GetCertificates()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(certs)
}

func TestClient_GetCertificate(t *testing.T) {
	c := Make(nil, &AuthHeader{
		Key:   "###",
		Value: "###",
	}, "###")
	cert, err := c.GetCertificate("###")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(cert)
}
