package lutron

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

func SendWait(conn *tls.Conn, r Request) (Response, error) {
	r.Header.ClientTag = newTag()

	b, err := requestJSON(r)
	if err != nil {
		return Response{}, err
	}

	if Debug {
		fmt.Println("sending request (raw):")
		fmt.Println(pretty(r))
	}

	if _, err = conn.Write(b); err != nil {
		return Response{}, fmt.Errorf("error writing to conn: %w", err)
	}

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(conn, &buf))

	rsp := Response{}

	// This is the wait part

	// TODO: Still need a timeout on the read... and the write above...

	deadline := time.Now().Add(RequestTimeout)
	for time.Now().Before(deadline) {
		if err := dec.Decode(&rsp); err != nil {
			return rsp, fmt.Errorf("error decoding response json: %w", err)
		}
		if Debug {
			fmt.Println("recieved response (raw):")
			fmt.Println(prettyb(buf.Bytes()))
			buf.Reset()
		}
		if rsp.Header.ClientTag == r.Header.ClientTag {
			if Debug {
				fmt.Println("response matches request ClientTag:", r.Header.ClientTag)
			}
			break
		}
	}

	return rsp, nil
}

func newTag() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func requestJSON(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request json: %w", err)
	}
	return append(b, '\r', '\n'), nil
}
