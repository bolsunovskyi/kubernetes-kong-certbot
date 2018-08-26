package kong

import "testing"

func TestClient_AddCertificate(t *testing.T) {
	c := Make(nil, nil, "http://###")
	if err := c.AddCertificate("foo", "bar", "###"); err != nil {
		t.Fatal(err)
	}
}
