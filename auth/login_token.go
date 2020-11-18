package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

const (
	issuer            = "pitbull"
	jwtPrivateKeyPath = "etc/key/jwt_rsa.pem"
	jwtPublicKeyPath  = "etc/key/jwt_rsa_pub.pem"
	jwtExpiry         = 10 * time.Minute
)

var (
	jwtSigningKey = readSigningKey()
	jwtPublicKey  = readPublicKey()

	contextKey = userClaimsContextKey("userclaims")
)

type userClaimsContextKey string

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func BuildJwt(username string) (string, error) {
	currentTime := time.Now()
	claims := UserClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: currentTime.Add(jwtExpiry).Unix(),
			IssuedAt:  currentTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	ss, err := token.SignedString(jwtSigningKey)
	if err != nil {
		log.WithError(err).Error("Failed to sign jwt")
		return "", err
	}

	return ss, nil
}

func VerifyJwt(token string) (*UserClaims, bool, error) {
	var claims UserClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtPublicKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, false, err
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, true, nil
			} else {
				return nil, false, err
			}
		} else {
			return nil, false, err
		}
	}

	return &claims, false, nil
}

func WithUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, contextKey, claims)
}

func ContextUserClaims(ctx context.Context) (*UserClaims, error) {
	if claims, ok := ctx.Value(contextKey).(*UserClaims); ok {
		return claims, nil
	}

	return nil, errors.New("no UserClaims in context")
}

func readSigningKey() *rsa.PrivateKey {
	raw := readKeyFile(jwtPrivateKeyPath)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(raw)
	if err != nil {
		log.WithError(err).Error("Unable to parse jwt signing private key")
		panic(err)
	}

	return key
}

func readPublicKey() *rsa.PublicKey {
	raw := readKeyFile(jwtPublicKeyPath)
	key, err := jwt.ParseRSAPublicKeyFromPEM(raw)
	if err != nil {
		log.WithError(err).Error("Unable to parse jwt public key")
		panic(err)
	}

	return key
}

func readKeyFile(path string) []byte {
	rawKey, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithError(err).Error("Unable to read jwt signing private key file")
		panic(err)
	}

	return rawKey
}
