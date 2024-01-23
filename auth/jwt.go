package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	UserId int `json:"userid"`
	jwt.StandardClaims
}

var (
	adminkey      = []byte("adminSecretKey")
	userKey       = []byte("userSecretKey")
	jwtvalidClaim *JwtClaim
	jwtclim       JwtClaim
)

func GenerateJwt(userId int, role string) (jwtToken string) {
	var key []byte
	expiryTime := time.Now().Add(1 * time.Hour)
	jwtclaim := JwtClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiryTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaim)
	if role == "admin" {
		key = adminkey
	} else {
		key = userKey
	}
	jwtToken, err := token.SignedString(key)
	if err != nil {
		return ""
	}
	return
}

func ValidateToken(jwtToken, role string) (int, error) {
	var (
		ok  bool
		key []byte
	)
	if role == "admin" {
		key = adminkey
	} else {
		key = userKey
	}
	jwtclim = JwtClaim{}
	token, err := jwt.ParseWithClaims(jwtToken, &jwtclim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return 0, err
	}

	if jwtvalidClaim, ok = token.Claims.(*JwtClaim); !ok {
		return 0, errors.New("Parsing error")
	}

	if jwtvalidClaim.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}

	return jwtvalidClaim.UserId, nil
}
