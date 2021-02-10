package useCase

import (
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"ExGabi/utils/token"
)

type UseCase struct{
	mutationRepository mutations.IRepository

}

func (uC *UseCase) AddItem(item *payload.Item) (*response.Item, error) {
	tokenClaims,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tokenClaims)
	if err!=nil{
		return nil, nil
	}
	id,err := uC.mutationRepository.AddItem(tokenClaims.Id,item)
	return &response.Item{Id: id, Token: newToken}, nil
}
func (uC *UseCase) DeleteItem(item *payload.Item) (string, error) {
	tokenClaims,err := token.CheckToken(item.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := token.CreateToken(tokenClaims)
	if err!=nil{
		return "", err
	}
	err =uC.mutationRepository.DeleteItem(tokenClaims.Id,item)
	return newToken,err
}
func (uC *UseCase) UpdateItem(item *payload.Item) (*response.Item, error) {
	tokenClaims,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tokenClaims)
	if err!=nil{
		return nil, err
	}
	newItem,err := uC.mutationRepository.UpdateItem(tokenClaims.Id,item)
	if err!=nil{
		return nil, err
	}
	newItem.Token=newToken
	return newItem,nil
}


func (uC *UseCase) DeleteUser(user *payload.User) (string, error) {
	tkn,err := token.CheckToken(user.Token)
	if err!=nil{
		return "",err
	}
	newToken,err := token.CreateToken(tkn)
	if err!=nil{
		return "", err
	}
	err = uC.mutationRepository.DeleteUser(user)
	return newToken,err
}
func (uC *UseCase) UpdateUser(user *payload.User) (*response.User, error) {
	tkn,err := token.CheckToken(user.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tkn)
	if err!=nil{
		return nil, err
	}
	newUser,err:= uC.mutationRepository.UpdateUser(user)
	if err!=nil{
		return nil, err
	}
	newUser.Token = newToken
	return newUser,nil
}


func (uC *UseCase) Register(user *payload.User) (string, error) {
	//err:= utils.CheckRegisterCredentials(user)
	//if err!=nil{
	//	return "", err
	//}
	//
	//id,err :=uC.mutationRepository.AddUser(user)
	//if err != nil{
	//	return "", err
	//}
	//responseUser,err := uC.mutationRepository.GetUserById(id)
	//if err!=nil{
	//	return "", nil
	//}
	//
	//tkn,err:=token.CreateToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	//if err!=nil{
	//	return "", nil
	//}
	return "",nil
}


func New(repo mutations.IRepository) mutations.IUseCase{
	return &UseCase{repo}
}

