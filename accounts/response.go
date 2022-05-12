package accounts

type errorResponse struct {
	Message string `json:"error_message,omitempty"`
}

type successResponse struct {
	Data interface{} `json:"data,omitempty"`
}
