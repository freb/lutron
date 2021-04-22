// doc comment
package lutron

import "time"

var (
	Debug = true

	KeyFile            = "caseta.key"
	CertFile           = "caseta.crt"
	CACertFile         = "caseta-bridge.crt"
	AppCN              = "lutron.go"
	LEAPPort           = "8081" // Lutron ? ? Protocol?
	LAPPort            = "8083" // Lutron Authentication Protocol? apk calls it Bridge Association Port
	RequestTimeout     = 10 * time.Second
	SocketTimeout      = 10 * time.Second
	ButtonPressTimeout = 180 * time.Second
)
