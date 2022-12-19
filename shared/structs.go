package shared

import "github.com/golang-jwt/jwt"

// TODO: replace all response payload with the appropriate struct below
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type LogPayload struct {
	Data   string
	Source string
}
type GatewayTokenJwtClaim struct {
	Gateway string
	jwt.StandardClaims
}

type APIResponse struct {
	Status  int      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
