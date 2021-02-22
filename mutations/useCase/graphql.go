package useCase

import (
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"ExGabi/utils"
	"ExGabi/utils/token"
)

type UseCase struct{
	mutationRepository mutations.IRepository

}

func New(repo mutations.IRepository) mutations.IUseCase{
	return &UseCase{repo}
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
	tknClaims,err := token.CheckToken(user.Token)
	if err!=nil{
		return nil, err
	}
	newToken,err := token.CreateToken(tknClaims)
	if err!=nil{
		return nil, err
	}
	user.Id = tknClaims.Id
	newUser,err:= uC.mutationRepository.UpdateUser(user)
	if err!=nil{
		return nil, err
	}
	newUser.Token = newToken
	return newUser,nil
}


func (uC *UseCase) Register(user *payload.User) (*response.User, error) {
	err:= utils.CheckRegisterCredentials(user)
	if err!=nil{
		return nil, err
	}

	id,err :=uC.mutationRepository.AddUser(user)
	if err != nil{
		return nil, err
	}
	responseUser,err := uC.mutationRepository.GetUserById(id)
	if err!=nil{
		return nil, err
	}

	tkn,err:=token.CreateToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	if err!=nil{
		return nil, err
	}
	return &response.User{Id: id,Token: tkn},nil
}
func (uC UseCase) Login(user *payload.User) (string, error) {
	responseUser,err := uC.mutationRepository.GetUserByCredentials(user)
	if err!=nil{
		return "", err
	}
	tkn,err:=token.CreateToken(&response.User{Id: responseUser.Id,Email: responseUser.Email})
	if err!=nil{
		return "", nil
	}
	return tkn,nil
}

func (uC *UseCase) GetMatchingSearch(item *payload.Item) (*[]response.Item,string, error) {
	tokenClaims,err := token.CheckToken(item.Token)
	if err!=nil{
		return nil,item.Token, err
	}
	newToken,err := token.CreateToken(tokenClaims)
	if err!=nil{
		return nil,item.Token,err
	}
	descriptionSearchCriteria := &payload.Item{Description: item.Title}

	descriptionItems,err := uC.mutationRepository.GetMatchingItems(tokenClaims.Id,descriptionSearchCriteria)

	if err!=nil{
		emptyList := new([]response.Item)
		return emptyList,newToken,err
	}


	return descriptionItems,newToken,nil
}

//func createIndexes(collectionName string, coll *mongo.Collection) error {
//	if collectionName == "quote_requests" {
//		mod1 := mongo.IndexModel{Keys: bson.M{"origin_name": 1}, Options: nil}
//		mod2 := mongo.IndexModel{Keys: bson.M{"destination_name": 1}, Options: nil}
//		_, err := coll.Indexes().CreateOne(context.TODO(), mod1)
//		if err != nil {
//			fmt.Println("fail create index for", collectionName)
//			return err
//		}
//		_, err = coll.Indexes().CreateOne(context.TODO(), mod2)
//		if err != nil {
//			fmt.Println("fail create index for", collectionName)
//			return err
//		}
//	}
//	return nil
//}