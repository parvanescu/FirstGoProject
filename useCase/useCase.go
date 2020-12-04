package useCase

import (
	"ExGabi/payload"
	"ExGabi/repository"
	"ExGabi/response"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type UseCase struct{
	itemRepository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)AddItem(item *payload.Item)(primitive.ObjectID,error){
	id,err := uC.itemRepository.AddItem(item)
	return id,err
}

func (uC *UseCase)DeleteItem(item *payload.Item)error {
	err :=uC.itemRepository.DeleteItem(item)
	return err
}

func (uC *UseCase)UpdateItem(item *payload.Item)(*response.Item,error) {
	newItem,err := uC.itemRepository.UpdateItem(item)
	return newItem,err
}

func(uC *UseCase)GetItemById(id primitive.ObjectID)(*response.Item,error) {
	list,err := uC.itemRepository.GetAllItems()
	if err!=nil{
		return &response.Item{}, err
	}
	for _,v := range *list{
		if v.Id == id{
			return &v,nil
		}
	}
	return &response.Item{}, errors.New("user not found")
}

func(uC *UseCase)GetAllItems()(*[]response.Item,error) {
	items,err := uC.itemRepository.GetAllItems()
	if err != nil{
		return nil, err
	}
	return items,nil
}

func(uC *UseCase) AddUser(user *payload.User)(primitive.ObjectID,error){
	id,err :=uC.itemRepository.AddUser(user)
	if err != nil{
		return [12]byte{}, err
	}
	return id,nil
}

func(uC *UseCase)DeleteUser(user *payload.User)error{
	err := uC.itemRepository.DeleteUser(user)
	return err
	}

func(uC *UseCase)UpdateUser(user *payload.User)(*response.User,error) {
	newUser,err:= uC.itemRepository.UpdateUser(user)
	return newUser,err
}

func(uC *UseCase)GetUserById(id primitive.ObjectID)(*response.User,error) {
	return uC.itemRepository.GetUserById(id)
	}

func(uC *UseCase)GetAllUsers()(*[]response.User,error) {
	users,err := uC.itemRepository.GetAllUsers()
	if err != nil{
		return nil,err
	}
	return users,nil
}

