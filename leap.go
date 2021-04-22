package lutron

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
)

// TODO: LAP and LEAP request and responses likely look different, if they are
// then separate them into two different objects.

type LEAPAuth struct {
	KeyPEM  []byte `json:"key_pem"`
	CertPEM []byte `json:"cert_pem"`
	RootPEM []byte `json:"root_pem"`
}

// LAPClient is a client for the LAP using for pairing and other
// pre-authentication requests.
type LEAPClient struct {
	addr  string
	conn  *tls.Conn
	certs LEAPAuth
}

func NewLEAPClient(addr string, certs LEAPAuth) *LEAPClient {
	return &LEAPClient{
		addr:  addr,
		certs: certs,
	}

}

func (c *LEAPClient) Connect() (*tls.Conn, error) {
	cert, err := tls.X509KeyPair(c.certs.CertPEM, c.certs.KeyPEM)
	if err != nil {
		return nil, fmt.Errorf("error loading client certs: %w", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(c.certs.RootPEM); !ok {
		return nil, errors.New("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", net.JoinHostPort(c.addr, LEAPPort), &tls.Config{
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

func (c *LEAPClient) SendWait(r Request) (Response, error) {
	if c.conn == nil {
		return Response{}, errors.New("not connected")
	}
	return SendWait(c.conn, r)
}

func (c *LEAPClient) Ping() error {
	req := Request{
		Header: Header{
			RequestType: Read,
			URL:         "/server/status/ping",
			// URL:         "/server/1/status/ping",
		},
		CommuniqueType: &ReadRequest,
		// Body:           Body{},
	}

	rsp, err := c.SendWait(req)
	if err != nil {
		return err
	}

	if rsp.Header.StatusCode != "200 OK" {
		return errors.New(rsp.Header.StatusCode)
	}

	// Has Body.PingResponse.LEAPVersion
	return nil
}

/*
URLS:
"CreateRequest", "/zone/{zone_id}/commandprocessor",
"CreateRequest", "/virtualbutton/{scene_id}/commandprocessor",
"ReadRequest", "/zone/{device['zone']}/status"
"ReadRequest", "/server/1/status/ping"
"ReadRequest", "/device"
"ReadRequest", "/server/2/id" // LIP device, pro only
"ReadRequest", "/area"
"ReadRequest", "/occupancygroup"
"ReadRequest", "/virtualbutton"
*/

// TODO: make command structs for each command type?
// Or possibly just make constructors for each that sets the CommandType and
// accepts the proper arguments

/*
Bodies
{
    "Command": {
        "CommandType": "GoToDimmedLevel",
        "DimmedLevelParameters": {
            "Level": value,
            "FadeTime": _format_duration(fade_time),
        },
    }
},

{
    "Command": {
        "CommandType": "GoToLevel",
        "Parameter": [{"Type": "Level", "Value": value}],
    }
},

{
    "Command": {
        "CommandType": "GoToFanSpeed",
        "FanSpeedParameters": {"FanSpeed": value},



{"Command": {"CommandType": "PressAndRelease"}},


*/

// def _format_duration(duration: timedelta) -> str:
//     """Convert a timedelta to the hh:mm:ss format used in LEAP."""
//     total_seconds = math.floor(duration.total_seconds())
//     seconds = int(total_seconds % 60)
//     total_minutes = math.floor(total_seconds / 60)
//     minutes = int(total_minutes % 60)
//     hours = int(total_minutes / 60)
//     return f"{hours:02}:{minutes:02}:{seconds:02}"
