package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

var key = []byte("akash")

func GenerateJwt(username string) (jwtToken string) {
	expiryTime := time.Now().Add( 1 * time.Hour)
	jwtclaim := JwtClaim{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiryTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaim)
	jwtToken, err := token.SignedString(key)
	if err != nil {
		fmt.Println("err :", err)
	}

	fmt.Println("jwtToken :", jwtToken)
	return
}

func ValidateToken(jwtToken string) error {

	var (
		jwtvalidClaim *JwtClaim
		ok            bool
	)

	jwtclaim := JwtClaim{}
	token, err := jwt.ParseWithClaims(jwtToken, &jwtclaim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		fmt.Println("Err is :", err)
		return err
	}

	if jwtvalidClaim, ok = token.Claims.(*JwtClaim); !ok {
		return errors.New("Parsing error")
	}

	if jwtvalidClaim.ExpiresAt < time.Now().Unix() {
		return errors.New("token expired")
	}
    fmt.Println("UserName :",jwtvalidClaim.UserName)
	return nil

}
