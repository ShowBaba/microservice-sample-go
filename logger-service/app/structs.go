package app

type LogPayload struct {
	Data   string `json:"data" validate:"required"`
	Source string `json:"source" validate:"required"`
}
