package enum

type TokenType string

const (
	TokenRefresh TokenType = "refresh-token"
	TokenReset   TokenType = "reset-token"
)
