package shared

// TODO: replace all response payload with the appropriate struct below
type JSONResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

type ListenerServicePayload struct {
	ServiceName string
	Data        interface{}
}