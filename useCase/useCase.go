package useCase

import (
	"ExGabi/payload"
	"ExGabi/repository"
	"ExGabi/response"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"strings"
	"time"
)


var jwtKey = []byte("ToDo_List_JWT_key")
type UseCase struct{
	itemRepository repository.IRepository

}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)AddItem(item *payload.Item)(*response.Item,error){
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return nil, nil
	}
	id,err := uC.itemRepository.AddItem(item)
	return &response.Item{Id: id, Token: newToken}, nil
}

func (uC *UseCase)DeleteItem(item *payload.Item)(string,error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return "", err
	}
	err =uC.itemRepository.DeleteItem(item)
	return newToken,err
}

func (uC *UseCase)UpdateItem(item *payload.Item)(*response.Item,error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.refreshToken(tkn)
	newItem,err := uC.itemRepository.UpdateItem(item)
	if err!=nil{
		return nil, err
	}
	newItem.Token=newToken
	return newItem,nil
}

func(uC *UseCase)GetItemById(item *payload.Item)(*response.Item,error) {
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return nil,err
	}
	getItem,err := uC.itemRepository.GetItemById(item.Id)
	if err!=nil{
		return nil,err
	}
	getItem.Token = newToken
	return getItem,nil
}

func(uC *UseCase)GetAllItems()(*[]response.Item,error) {
	items,err := uC.itemRepository.GetAllItems()
	if err != nil{
		return nil, err
	}
	return items,nil
}



func(uC *UseCase)DeleteUser(user *payload.User)(string,error){
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return "", err
	}
	err = uC.itemRepository.DeleteUser(user)
	return newToken,err
	}

func(uC *UseCase)UpdateUser(user *payload.User)(*response.User,error) {
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return nil, err
	}
	newUser,err:= uC.itemRepository.UpdateUser(user)
	if err!=nil{
		return nil, err
	}
	newUser.Token = newToken
	return newUser,nil
}

func(uC *UseCase)GetUserById(user *payload.User)(*response.User,error) {
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := uC.refreshToken(tkn)
	if err!=nil{
		return nil, err
	}
	usr,err:=uC.itemRepository.GetUserById(user.Id)
	if err!=nil{
		return nil,err
	}
	usr.Token=newToken
	return usr,nil
	}

func(uC *UseCase)GetAllUsers()(*[]response.User,error) {
	users,err := uC.itemRepository.GetAllUsers()
	if err != nil{
		return nil,err
	}
	return users,nil
}


func(uC *UseCase) Register(user *payload.User)(string,error){
	err:= uC.checkRegisterCredentials(user)
	if err!=nil{
		return "", err
	}

	id,err :=uC.itemRepository.AddUser(user)
	if err != nil{
		return "", err
	}
	responseUser,err := uC.itemRepository.GetUserById(id)
	if err!=nil{
		return "", nil
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	responseUser.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,responseUser)
	tokenString,err := token.SignedString(jwtKey)
	if err!=nil{
		return "", nil
	}
	return tokenString,nil
}
func(uC * UseCase) Login(user *payload.User)(string,error){

	responseUser,err := uC.itemRepository.GetUserByCredentials(user)
	if err!=nil{
		return "", err
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	responseUser.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,responseUser)
	tokenString,err := token.SignedString(jwtKey)
	if err!=nil{
		return "", nil
	}
	return tokenString,nil

}

func(uC *UseCase)checkRegisterCredentials(user * payload.User)error{
	if len(user.FirstName) < 3{
		return errors.New("first name's length is too small")
	}
	if len(user.LastName) < 3{
		return errors.New("last name's length is to small")
	}
	if strings.Contains(user.Email,"@") == false{
		return errors.New("invalid email address")
	}
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+=?^_`{|}~-]+@[a-zA-Z0-9.]*.[a-z]+$")
	if !emailRegex.MatchString(user.Email){
		return errors.New("invalid email address")
	}
	if len(user.Password) < 6{
		return errors.New("password's length is to small")
	}
	return nil
}
func(uC *UseCase)checkToken(tkn string) (*response.User,error) {
	claims := &response.User{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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
func(uC *UseCase)refreshToken(user *response.User)(string,error){

	expirationTime := time.Now().Add(5 * time.Minute)
	user.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,user)
	tokenString,err := token.SignedString(jwtKey)
	if err!=nil{
		return "", nil
	}
	return tokenString,nil
}
