package collector

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	pb "github.com/icphalanx/rpc"
)

type Collector struct {
	conf     Config
	ingestor Ingestor
}

type CfsslSignResponse struct {
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Result   *struct {
		Certificate string `json:"certificate"`
	} `json:"result"`
	Success bool `json:"success"`
}

func NewServer(conf Config) (*Collector, error) {
	c := new(Collector)
	c.conf = conf
	c.ingestor = Ingestor{
		c.conf.DBConn,
	}
	c.setupWeb()

	return c, c.addOwnSSLCertificate()
}

func (c *Collector) addOwnSSLCertificate() error {
	if c.conf.TLSConfig.Certificates[0].Leaf == nil {
		var err error
		c.conf.TLSConfig.Certificates[0].Leaf, err = x509.ParseCertificate(c.conf.TLSConfig.Certificates[0].Certificate[0])
		if err != nil {
			return err
		}
	}
	_, err := c.ingestor.GetOrAddSSLCertificate(c.conf.TLSConfig.Certificates[0].Leaf)
	return err
}

var (
	ErrBadAuth                  = errors.New(`non-permitted AuthType detected`)
	ErrProvisioningCert         = errors.New(`provisioning certificates are not valid for this operation`)
	ErrProvisioningCertRequired = errors.New(`you must use a provisioning certifiate for this operation`)
)

const (
	ProvisioningCertCommonName = `Phalanx Provisioning Certificate`
)

func getTLSInfo(ctx context.Context) (*credentials.TLSInfo, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrBadAuth
	}

	if p.AuthInfo.AuthType() != "tls" {
		return nil, ErrBadAuth
	}

	t := p.AuthInfo.(credentials.TLSInfo)
	return &t, nil
}

func getCertificateFromContext(ctx context.Context) (*x509.Certificate, error) {
	t, err := getTLSInfo(ctx)
	if err != nil {
		return nil, err
	}

	return t.State.PeerCertificates[0], nil
}

func validateStreamIsFullCert(ctx context.Context) error {
	cert, err := getCertificateFromContext(ctx)
	if err != nil {
		return err
	}

	if cert.Subject.CommonName == ProvisioningCertCommonName {
		return ErrProvisioningCert
	}
	return nil
}

func validateStreamIsProvisioningCertOrExpiring(ctx context.Context) error {
	cert, err := getCertificateFromContext(ctx)
	if err != nil {
		return err
	}

	if cert.Subject.CommonName == ProvisioningCertCommonName {
		return nil
	}

	if cert.NotAfter.Before(time.Now().Add(35 * 24 * time.Hour)) {
		return nil
	}

	return ErrProvisioningCertRequired
}

func getRemoteHostnameFromCert(ctx context.Context) (string, error) {
	cert, err := getCertificateFromContext(ctx)
	if err != nil {
		return "", err
	}

	cn := cert.Subject.CommonName
	if cn == ProvisioningCertCommonName {
		return cn, ErrProvisioningCert
	}

	return cn, nil
}

func pemCertToX509(pemCert string) (*x509.Certificate, error) {
	pemBlock, _ := pem.Decode([]byte(pemCert))
	if pemBlock == nil {
		return nil, fmt.Errorf("PEM failed to decode")
	}
	if pemBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("PEM block type was %s not 'CERTIFICATE'", pemBlock.Type)
	}

	return x509.ParseCertificate(pemBlock.Bytes)
}

func (c *Collector) ensureSSLCertificateFromContext(ctx context.Context) error {
	certFromCtx, err := getCertificateFromContext(ctx)
	if err != nil {
		return err
	}

	dbcert, err := c.ingestor.GetOrAddSSLCertificate(certFromCtx)
	if err != nil {
		return err
	}

	if dbcert.Revoked {
		return fmt.Errorf("peer certificate marked as revoked")
	}

	return err
}

func (c *Collector) ConfigureMe(ctx context.Context, host *pb.Host) (*pb.HostConfiguration, error) {
	if err := validateStreamIsFullCert(ctx); err != nil {
		return nil, err
	}

	remoteHost, err := getRemoteHostnameFromCert(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.ensureSSLCertificateFromContext(ctx); err != nil {
		return nil, err
	}

	log.Printf("%s: ConfigureMe{%s}", remoteHost, host.String())

	hc := new(pb.HostConfiguration)
	hc.Port = 0
	return hc, nil
}

func (c *Collector) Report(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	if err := validateStreamIsFullCert(ctx); err != nil {
		return nil, err
	}

	remoteHost, err := getRemoteHostnameFromCert(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.ensureSSLCertificateFromContext(ctx); err != nil {
		return nil, err
	}

	log.Printf("%s: Report{%s}", remoteHost, req.String())

	err = c.ingestor.IngestReport(req)
	if err != nil {
		log.Printf("%s [Report{%s}]: error: %v", remoteHost, req.String(), err)
	}

	resp := new(pb.ReportResponse)
	resp.Success = err == nil
	return resp, nil
}

func (c *Collector) SignMe(ctx context.Context, req *pb.SigningRequest) (*pb.SigningResponse, error) {
	if err := validateStreamIsProvisioningCertOrExpiring(ctx); err != nil {
		log.Println("Got SignMe request from invalid certificate")
		return nil, err
	}
	log.Println("Got SignMe request with certificate", req.Csr)

	jsonReq := map[string]string{
		"certificate_request": req.Csr,
	}
	reqBytes, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, err
	}

	// contact the CFSSL server!
	b := bytes.NewReader(reqBytes)
	resp, err := http.Post(fmt.Sprintf("%s/api/v1/cfssl/sign", c.conf.CARemote), "application/json", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jresp := CfsslSignResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&jresp); err != nil {
		return nil, err
	}
	log.Printf("Got response from cfssl %#v", jresp)

	if !jresp.Success {
		if len(jresp.Errors) > 0 {
			return nil, fmt.Errorf("%v", jresp.Errors[0])
		}
		return nil, fmt.Errorf("something was wrong with the request to cfssl")
	}

	signingResponse := new(pb.SigningResponse)
	signingResponse.Cert = jresp.Result.Certificate

	// reokve the old cert
	certFromCtx, err := getCertificateFromContext(ctx)
	if err != nil {
		log.Printf("got error %v whilst trying to get certificate from ctx", err)
		return signingResponse, nil
	}

	err = c.ingestor.RevokeCertificate(certFromCtx)
	if err != nil {
		log.Printf("got error %v whilst trying to revoke old certificate :(", err)
	}

	// we should add this certificate to the database and revoke the old one
	cert, err := pemCertToX509(jresp.Result.Certificate)
	if err != nil {
		// we'll probably catch it next time
		log.Printf("got error %v whilst trying to parse certificate from cfssl to add to DB")
		return signingResponse, nil
	}

	_, err = c.ingestor.GetOrAddSSLCertificate(cert)
	if err != nil {
		log.Printf("got error %v whilst trying to add certificate to DB from cfssl", err)
		return signingResponse, nil
	}

	return signingResponse, nil
}

func (c *Collector) RecordLogs(rls pb.PhalanxCollector_RecordLogsServer) error {
	for {
		ln, err := rls.Recv()
		if err != nil {
			return err
		}

		err = c.ingestor.IngestLine(ln)
		if err != nil {
			return err
		}

		MainHub.Send(*ln, []string{ln.Host})
	}
}
