package config

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTToken(email string) (string, errs.AppError) {
	secretApiKey, err_ := Value(key.SecretApiKey)
	if err_ != nil {
		err_ = err_.Annotatef(applog.Log, "unable to get env variables config: %v", err_)
		return "", err_
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().AddDate(0, 0, 10)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, _err := token.SignedString([]byte(secretApiKey))
	if _err != nil {
		// If there is an error in creating the JWT return an internal server error
		err := errs.ErrCreatingJWT.Throwf(applog.Log, errs.ErrFmt, _err.Error())
		return "", err
	}

	return tokenString, nil
}
