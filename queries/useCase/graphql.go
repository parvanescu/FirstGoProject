package useCase

import (
	"ExGabi/payload"
	"ExGabi/queries"
	"ExGabi/response"
	"ExGabi/utils/token"
)

type UseCase struct{
	queriesRepository queries.IRepository

}

func (uC *UseCase) GetItemById(item *payload.Item) (*response.Item, error) {
	tkn,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tkn)
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
	tkn,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tkn)
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
	tkn,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil,item.Token,err
	}
	newToken,err := token.CreateToken(tkn)
	if err!=nil{
		return nil,item.Token,err
	}
	getItem,err := uC.queriesRepository.GetItemByDescription(item)
	if err!=nil{
		return nil,newToken,err
	}
	return getItem,newToken,nil
}
func (uC *UseCase) GetAllItems(payloadToken string) (*[]response.Item, string, error) {
	tkn,err := token.CheckToken(payloadToken)
	if err!=nil {
		return nil,"",err
	}
	newToken,err:=token.CreateToken(tkn)
	if err != nil{
		return nil,newToken,err
	}

	items,err := uC.queriesRepository.GetAllItems()

	return items,newToken,nil
}
func (uC *UseCase) GetAllUsersItems(payloadToken string) (*[]response.Item,string,error){
	tknClaims,err := token.CheckToken(payloadToken)
	if err!=nil {
		return nil,"",err
	}
	newToken,err:=token.CreateToken(tknClaims)
	if err != nil{
		return nil,newToken,err
	}

	items,err := uC.queriesRepository.GetAllUsersItems(tknClaims.Id)

	return items,newToken,nil
}

func (uC *UseCase) GetUserById(user *payload.User) (*response.User, error) {
	tkn,err := token.CheckToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := token.CreateToken(tkn)
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
func (uC *UseCase)GetUserProfile(payloadToken string)(*response.User,error){
	tknClaims,err := token.CheckToken(payloadToken)
	if err!=nil{
		return nil,err
	}
	newToken,err := token.CreateToken(tknClaims)
	if err!=nil{
		return nil, err
	}
	usr,err:=uC.queriesRepository.GetUserById(tknClaims.Id)
	if err!=nil{
		return nil,err
	}
	usr.Token=newToken
	return usr,nil
}
func (uC *UseCase) GetUserByEmail(user *payload.User) (*response.User, error) {
	tkn,err := token.CheckToken(user.Token)
	if err!=nil{
		return nil,err
	}
	newToken,err := token.CreateToken(tkn)
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
func (uC *UseCase) GetAllUsers(payloadToken string) (*[]response.User, string, error) {
	tkn,err := token.CheckToken(payloadToken)
	if err!=nil {
		return nil,"",err
	}

	newToken,err:=token.CreateToken(tkn)
	if err!=nil{
		return nil, newToken, err
	}

	users,err := uC.queriesRepository.GetAllUsers()
	if err != nil{
		return nil,newToken,err
	}

	return users,newToken,nil
}


func New(repo queries.IRepository)queries.IUseCase{
	return &UseCase{repo}
}

