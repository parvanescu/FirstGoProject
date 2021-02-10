package token

import (
	"ExGabi/response"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)
var JwtKey = []byte("ToDo_List_JWT_key")

func CheckToken(tkn string) (*response.User,error) {
	claims := &response.User{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil,errors.New("invalid signature")
		}
		return nil,err
	}
	if !token.Valid {
		return nil,errors.New("invalid token")
	}
	return claims,nil
}
func CreateToken(user *response.User)(string,error) {

	expirationTime := time.Now().Add(5 * time.Hour)
	user.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}