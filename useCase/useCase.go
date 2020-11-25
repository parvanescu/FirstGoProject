package useCase

import (
	"ExGabi/payload"
	"ExGabi/repository"
	"ExGabi/response"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type UseCase struct{
	itemRepository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)AddItem(item payload.Item){
	uC.itemRepository.AddItem(item)
}

func (uC *UseCase)DeleteItem(userId primitive.ObjectID,itemId primitive.ObjectID)error {
	err :=uC.itemRepository.DeleteItem(userId,itemId)
	return err
}

func (uC *UseCase)UpdateItem(id primitive.ObjectID,item payload.Item)(response.Item,error) {
	newItem,err := uC.itemRepository.UpdateItem(id,item)
	return newItem,err
}

func(uC *UseCase)GetItem(id primitive.ObjectID)response.Item {
	
	for _,v := range *uC.itemRepository.GetAllItems(){
		if v.ItemId == id{
			return v
		}
	}
	return response.NewItem("","")
}

func(uC *UseCase)GetAllItems() {
	items := *uC.itemRepository.GetAllItems()
	for _, v := range items {
		fmt.Println(v.ItemId.Hex())
		fmt.Println(v.Title)
		fmt.Println(v.Description)
	}
}

func(uC *UseCase) AddUser(user payload.User){
	uC.itemRepository.AddUser(user) //returns object id of user
	}

func(uC *UseCase)DeleteUser(id primitive.ObjectID)error{
	var err error
	user := uC.itemRepository.GetUserById(id)
	for _,v := range user.Items{
		err =uC.DeleteItem(id,v.ItemId)
	}
	if err != nil{
		return err
	}
	err = uC.itemRepository.DeleteUser(id)
	return err
	}

func(uC *UseCase)UpdateUser(id primitive.ObjectID,user payload.User)(response.User,error) {
	newUser,err:= uC.itemRepository.UpdateUser(id,user)
	return newUser,err
}

func(uC *UseCase)GetUser(id primitive.ObjectID)response.User {
	return uC.itemRepository.GetUserById(id)
	}

func(uC *UseCase)GetAllUsers() {
	users := *uC.itemRepository.GetAllUsers()
	for _,v := range users  {
		fmt.Println(v.UserId.Hex())
		fmt.Println(v.UserName)
		fmt.Println(v.Items)
	}
	}

