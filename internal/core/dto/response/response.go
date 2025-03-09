package response

type Response struct {
	Status  string            `json:"status"`
	Errors  map[string]string `json:"errors,omitempty"`
	Message string            `json:"message,omitempty"`
}
