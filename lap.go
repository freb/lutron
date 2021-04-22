package lutron

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// LAPClient is a client for the LAP using for pairing and other
// pre-authentication requests.
type LAPClient struct {
	addr string
	conn *tls.Conn
}

func NewLAPClient(addr string) *LAPClient {
	return &LAPClient{addr: addr}

}

func (c *LAPClient) Connect() (*tls.Conn, error) {
	cert, err := tls.X509KeyPair(LAPCertPEM, LAPKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("error loading app certs: %w", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(CasetaRootCertPEM); !ok {
		return nil, errors.New("failed to parse root certificate")
	}

	// TODO: set dealines on connection, then use: tls.DialWithDialer()

	conn, err := tls.Dial("tcp", net.JoinHostPort(c.addr, LAPPort), &tls.Config{
		RootCAs:      roots,
		Certificates: []tls.Certificate{cert},
		// We cannot verify hostname because server cert CN's are not for a specific hostname.
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, fmt.Errorf("error connecting to bridge: %w", err)
	}

	c.conn = conn

	return conn, nil
}

func (c *LAPClient) SendWait(r Request) (Response, error) {
	if c.conn == nil {
		return Response{}, errors.New("not connected")
	}
	return SendWait(c.conn, r)
}

var csrTemplate = &x509.CertificateRequest{
	SignatureAlgorithm: x509.SHA256WithRSA, // need?
	Subject:            pkix.Name{CommonName: AppCN},
	// APK uses these details:
	// COMMON_NAME = "Lutron App";
	// COUNTRY = "US";
	// LOCATION = "Coopersburg";
	// ORG_NAME = "Lutron Electronics Co.\\, Inc.";
	// STATE = "Pennsylvania";
}

func (c *LAPClient) Pair() (LEAPAuth, error) {
	certs := LEAPAuth{}
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return certs, fmt.Errorf("error generating private key: %w", err)
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, pk)
	if err != nil {
		return certs, fmt.Errorf("error generating csr: %w", err)
	}

	// When to do connect, do it in here?

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})

	fmt.Println("Press the small black button on the back of the Caseta bridge...")

	// After connecting, once the button is pressed the server will send data.
	// Since there is no initial request, this is presumably why they use a
	// spearate port.

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(c.conn, &buf))

	rsp := Response{}

	// TODO: for t < timeout, since we're looping responses we cannot simply
	// use the socket timeout
	deadline := time.Now().Add(RequestTimeout)
	for time.Now().Before(deadline) {
		if err := dec.Decode(&rsp); err != nil {
			return certs, fmt.Errorf("error decoding response json: %w", err)
		}

		if Debug {
			fmt.Println("recieved response (raw):")
			fmt.Println(prettyb(buf.Bytes()))
			buf.Reset()
		}

		if !strings.HasPrefix(rsp.Header.ContentType, "status;") {
			continue
		}

		found := false
		for _, perm := range rsp.Body.Status.Permissions {
			if perm == "PhysicalAccess" {
				found = true
			}
		}
		if found {
			break
		}
	}

	req := Request{
		Header: Header{
			RequestType: Execute,
			URL:         "/pair",
			ClientTag:   newTag(),
		},
		Body: &Body{
			CommandType: CSR,
			Parameters: map[string]interface{}{
				"CSR":         string(csrPEM),
				"DisplayName": "lutron.go",
				"DeviceUID":   "000000000000",
				"Role":        "Admin",
			},
		},
	}

	rsp, err = c.SendWait(req)
	if err != nil {
		return certs, err
	}

	certs.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	certs.CertPEM = []byte(rsp.Body.SigningResult.Certificate)
	certs.RootPEM = []byte(rsp.Body.SigningResult.RootCertificate)

	return certs, nil
}

func (c *LAPClient) FetchBridgeCert() (Response, error) {
	req := Request{
		Header: Header{
			RequestType: Read,
			URL:         "/certificate/root",
		},
	}

	return c.SendWait(req)
}

func (c *LAPClient) GetMacAddress() (Response, error) {
	req := Request{
		Header: Header{
			RequestType: Read,
			URL:         "/system/macaddress",
		},
	}

	return c.SendWait(req)
}

func (c *LAPClient) GetCrossSign() (Response, error) {
	req := Request{
		Header: Header{
			RequestType: Read,
			URL:         "/system/status/crosssign",
		},
	}

	return c.SendWait(req)
}
