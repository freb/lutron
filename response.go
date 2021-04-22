package lutron

type Response struct {
	Header struct {
		StatusCode  string `json:"StatusCode,omitempty"`
		ContentType string `json:"ContentType,omitempty"`
		ClientTag   string `json:"ClientTag,omitempty"`
	} `json:"Header,omitempty"`
	CommuniqueType string `json:"CommuniqueType,omitempty"`
	Body           struct {
		Status struct {
			Permissions []string `json:"Permissions,omitempty"`
		} `json:"Status,omitempty"`
		Exception struct {
			Message string `json:"Message,omitempty"`
		} `json:"Exception,omitempty"`
		SigningResult struct {
			Certificate     string `json:"Certificate,omitempty"`
			RootCertificate string `json:"RootCertificate,omitempty"`
		} `json:"SigningResult,omitempty"`
		PingResponse struct {
			LEAPVersion float32 `json:"LEAPVersion"`
		} `json:"PingResponse,omitempty"`
	} `json:"Body,omitempty"`
}
