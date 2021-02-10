package useCase

import (
	"ExGabi/payload"
	"ExGabi/queries"
	"ExGabi/response"
	token2 "ExGabi/utils/token"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UseCase struct{
	queriesRepository queries.IRepository

}

func (uC *UseCase) GetItemById(item *payload.Item) (*response.Item, error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil,err
	}
	getItem,err := uC.queriesRepository.GetItemById(item.Id)
	if err!=nil{
		return nil,err
	}
	getItem.Token = newToken
	return getItem,nil
}
func (uC *UseCase) GetItemByTitle(item *payload.Item) (*response.Item, error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil,err
	}
	getItem,err := uC.queriesRepository.GetItemByTitle(item)
	if err!=nil{
		return nil,err
	}
	getItem.Token = newToken
	return getItem,nil
}
func (uC *UseCase) GetItemByDescription(item *payload.Item) (*[]response.Item, string, error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil,item.Token,err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil,item.Token,err
	}
	getItem,err := uC.queriesRepository.GetItemByDescription(item)
	if err!=nil{
		return nil,newToken,err
	}
	return getItem,newToken,nil
}
func (uC *UseCase) GetAllItems(token string) (*[]response.Item, string, error) {
	tkn,err := uC.checkToken(token)
	if err!=nil {
		return nil,"",err
	}
	newToken,err:=uC.createToken(tkn)
	if err != nil{
		return nil,newToken,err
	}

	items,err := uC.queriesRepository.GetAllItems()

	return items,newToken,nil
}


func (uC *UseCase) GetUserById(user *payload.User) (*response.User, error) {
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil, err
	}
	usr,err:=uC.queriesRepository.GetUserById(user.Id)
	if err!=nil{
		return nil,err
	}
	usr.Token=newToken
	return usr,nil
}
func (uC *UseCase) GetUserByEmail(user *payload.User) (*response.User, error) {
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil, err
	}
	usr,err:=uC.queriesRepository.GetUserByEmail(user)
	if err!=nil{
		return nil,err
	}
	usr.Token=newToken
	return usr,nil
}
func (uC *UseCase) GetAllUsers(token string) (*[]response.User, string, error) {
	tkn,err := uC.checkToken(token)
	if err!=nil {
		return nil,"",err
	}

	newToken,err:=uC.createToken(tkn)
	if err!=nil{
		return nil, newToken, err
	}

	users,err := uC.queriesRepository.GetAllUsers()
	if err != nil{
		return nil,newToken,err
	}

	return users,newToken,nil
}


func (uC UseCase) Login(user *payload.User) (string, error) {
	responseUser,err := uC.queriesRepository.GetUserByCredentials(user)
	if err!=nil{
		return "", err
	}
	tkn,err:=uC.createToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	if err!=nil{
		return "", nil
	}
	return tkn,nil
}

func New(repo queries.IRepository)queries.IUseCase{
	return &UseCase{repo}
}

func(uC *UseCase)checkToken(tkn string) (*response.User,error) {
	claims := &response.User{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
		return token2.JwtKey, nil
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
func(uC *UseCase)createToken(user *response.User)(string,error){

	expirationTime := time.Now().Add(5 * time.Hour)
	user.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,user)
	tokenString,err := token.SignedString(token2.JwtKey)
	if err!=nil{
		return "", nil
	}
	return tokenString,nil
}
