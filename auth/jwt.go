package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	UserId int `json:"userid"`
	jwt.StandardClaims
}

var (
	adminkey           = []byte("adminSecretKey")
	userKey            = []byte("userSecretKey")
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
	if role == "admin"{
		key = adminkey
	}else{
		key = userKey
	}
	jwtToken, err := token.SignedString(key)
	if err != nil {
		fmt.Println("err :", err)
	}

	fmt.Println("jwtToken :", jwtToken)
	return
}

func ValidateToken(jwtToken, role string) error {

	var (
		ok bool
		key []byte
	)
    if role == "admin"{
		key = adminkey
	}else{
		key = userKey
	}
	jwtclim = JwtClaim{}
	token, err := jwt.ParseWithClaims(jwtToken, &jwtclim, func(t *jwt.Token) (interface{}, error) {
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

	return nil

}

func GetUserId() int {
	return jwtclim.UserId
}
