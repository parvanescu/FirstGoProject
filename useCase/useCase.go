package useCase

import (
	"ExGabi/payload"
	"ExGabi/repository"
	"ExGabi/response"
	"ExGabi/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
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
	tokenClaims,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.createToken(tokenClaims)
	if err!=nil{
		return nil, nil
	}
	id,err := uC.itemRepository.AddItem(tokenClaims.Id,item)
	return &response.Item{Id: id, Token: newToken}, nil
}
func (uC *UseCase)DeleteItem(item *payload.Item)(string,error) {
	tokenClaims,err := uC.checkToken(item.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := uC.createToken(tokenClaims)
	if err!=nil{
		return "", err
	}
	err =uC.itemRepository.DeleteItem(tokenClaims.Id,item)
	return newToken,err
}
func (uC *UseCase)UpdateItem(item *payload.Item)(*response.Item,error) {
	tokenClaims,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.createToken(tokenClaims)

	_,err = uC.itemRepository.GetItemById(item.Id)
	if err !=nil{
		return &response.Item{Token: newToken},err
	}

	newItem,err := uC.itemRepository.UpdateItem(tokenClaims.Id,item)
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
	newToken,err := uC.createToken(tkn)
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
func(uC *UseCase)GetItemByTitle(item *payload.Item)(*response.Item,error){
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil,err
	}
	getItem,err := uC.itemRepository.GetItemByTitle(item)
	if err!=nil{
		return nil,err
	}
	getItem.Token = newToken
	return getItem,nil
}
func(uC *UseCase)GetItemByDescription(item *payload.Item)(*[]response.Item,string,error){
	tkn,err := uC.checkToken(item.Token)
	if err!=nil{
		return nil,item.Token,err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil,item.Token,err
	}
	getItem,err := uC.itemRepository.GetItemByDescription(item)
	if err!=nil{
		return nil,newToken,err
	}
	return getItem,newToken,nil
}
func(uC *UseCase)GetAllItems(token string)(*[]response.Item,string,error) {
	tkn,err := uC.checkToken(token)
	if err!=nil {
		return nil,"",err
	}
	newToken,err:=uC.createToken(tkn)
	if err != nil{
		return nil,newToken,err
	}

	items,err := uC.itemRepository.GetAllItems()

	return items,newToken,nil
}



func(uC *UseCase)DeleteUser(user *payload.User)(string,error){
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := uC.createToken(tkn)
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
	newToken,err := uC.createToken(tkn)
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
	newToken,err := uC.createToken(tkn)
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
func(uC *UseCase)GetUserByEmail(user *payload.User)(*response.User,error){
	tkn,err := uC.checkToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := uC.createToken(tkn)
	if err!=nil{
		return nil, err
	}
	usr,err:=uC.itemRepository.GetUserByEmail(user)
	if err!=nil{
		return nil,err
	}
	usr.Token=newToken
	return usr,nil
}
func(uC *UseCase)GetAllUsers(token string)(*[]response.User,string,error) {
	tkn,err := uC.checkToken(token)
	if err!=nil {
		return nil,"",err
	}

	newToken,err:=uC.createToken(tkn)
	if err!=nil{
		return nil, newToken, err
	}

	users,err := uC.itemRepository.GetAllUsers()
	if err != nil{
		return nil,newToken,err
	}

	return users,newToken,nil
}


func(uC *UseCase) Register(user *payload.User)(string,error){
	err:= utils.CheckRegisterCredentials(user)
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

	tkn,err:=uC.createToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	if err!=nil{
		return "", nil
	}
	return tkn,nil
}
func(uC * UseCase) Login(user *payload.User)(string,error){

	responseUser,err := uC.itemRepository.GetUserByCredentials(user)
	if err!=nil{
		return "", err
	}
	tkn,err:=uC.createToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	if err!=nil{
		return "", nil
	}
	return tkn,nil
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
func(uC *UseCase)createToken(user *response.User)(string,error){

	expirationTime := time.Now().Add(5 * time.Hour)
	user.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,user)
	tokenString,err := token.SignedString(jwtKey)
	if err!=nil{
		return "", nil
	}
	return tokenString,nil
}
