package collector

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func generateCertPoolFromPath(caPath string) (*x509.CertPool, error) {
	b, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("failed to load certificate from PEM")
	}
	return cp, nil
}

func generateTLSConfig(caCertPool *x509.CertPool, cert tls.Certificate) (*tls.Config, error) {
	c := new(tls.Config)
	c.RootCAs = caCertPool
	c.Certificates = []tls.Certificate{cert}

	c.ClientCAs = caCertPool
	c.ClientAuth = tls.RequireAndVerifyClientCert
	return c, nil
}

func MakeTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	caPool, err := generateCertPoolFromPath(caPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return generateTLSConfig(caPool, cert)
}
