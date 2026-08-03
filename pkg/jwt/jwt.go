package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Parser struct {
	secretKey []byte
}

func NewParser(secretKey string) *Parser {
	return &Parser{secretKey: []byte(secretKey)}
}

func (p *Parser) Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return p.secretKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, fmt.Errorf("jwt.Parse: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (p *Parser) NewToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}
