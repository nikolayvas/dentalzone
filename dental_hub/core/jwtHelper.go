package core

import (
	"net/http"

	config "dental_hub/configuration"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var (
	JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetInstance().JwtSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
)

// GenerateJwt generates new JWT token
func GenerateJwt(dentistID string) (*string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"dentistId": dentistID})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.GetInstance().JwtSecret))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// ExtractJwtClaim extracts claim from authentication token
func ExtractJwtClaim(r *http.Request, claim string) (*string, error) {
	tokenString, err := JwtMiddleware.Options.Extractor(r)

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, JwtMiddleware.Options.ValidationKeyGetter)

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	p := claims[claim].(string)
	return &p, nil
}
