package lutron

type CommuniqueType string

var (
	CreateRequest    CommuniqueType = "CreateRequest"
	ReadRequest      CommuniqueType = "ReadRequest"
	SubscribeRequest CommuniqueType = "SubscribeRequest"
)

type Request struct {
	Header         Header          `json:"Header,omitempty"`
	CommuniqueType *CommuniqueType `json:"CommuniqueType,omitempty"`
	Body           *Body           `json:"Body,omitempty"`
}

type RequestType string

var (
	Read      RequestType = "Read"
	Execute   RequestType = "Execute"
	Subscribe RequestType = "Subscribe"
)

type Header struct {
	RequestType RequestType `json:"RequestType"`
	URL         string      `json:"Url"`
	ClientTag   string      `json:"ClientTag"` // generated by client to identify response
}

type Body struct {
	Command     Command                `json:"Command"`
	CommandType CommandType            `json:"CommandType"` // LAP
	Parameters  map[string]interface{} `json:"Parameters"`  // LAP
}

type CommandType string

var (
	CSR             CommandType = "CSR"
	GoToDimmedLevel CommandType = "GoToDimmedLevel"
	GoToFanSpeed    CommandType = "GoToFanSpeed"
	GoToLevel       CommandType = "GoToLevel"
	PressAndRelease CommandType = "PressAndRelease"
)

type DimmedLevelParameters struct {
	Level    int    `json:"Level"`
	FadeTime string `json:"FadeTime"`
}

type FanSpeedParameters struct {
	FanSpeed string `json:"FanSpeed"`
}

type Parameter struct {
	Type  string `json:"Parameter"`
	Value int    `json:"Value"`
}

type Command struct {
	CommandType           CommandType           `json:"CommandType"`
	DimmedLevelParameters DimmedLevelParameters `json:"DimmedLevelParameters"`
	FanSpeedParameters    FanSpeedParameters    `json:"FanSpeedParameters"`
	Parameter             []Parameter           `json:"Parameter"`
}
