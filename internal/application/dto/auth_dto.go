package dto

import "github.com/golang-jwt/jwt/v5"

type RESTClaims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
