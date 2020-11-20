package useCase

import (
	"ExGabi/model"
	"ExGabi/repository"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUseCase interface {
	AddItem(title string,description string)
	DeleteItem(id primitive.ObjectID)model.Item
	UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)model.Item
	GetItem(id primitive.ObjectID)model.Item
	GetAll() //[]model.Item
}

type UseCase struct{
	itemRepository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)AddItem(title string,description string){
	var item model.Item = model.NewItem(title,description)
	uC.itemRepository.AddItem(item)
}

func (uC *UseCase)DeleteItem(id primitive.ObjectID)model.Item {
	return uC.itemRepository.DeleteItem(id)
}

func (uC *UseCase)UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)model.Item {
	return uC.itemRepository.UpdateItem(id,model.NewItem(newTitle,newDescription))
}

func(uC *UseCase)GetItem(id primitive.ObjectID)model.Item {
	
	for _,v := range uC.itemRepository.GetAllItems(){
		if v.ItemId == id{
			return v
		}
	}
	return model.NewItem("","")
}

func(uC *UseCase)GetAll(){
	items := uC.itemRepository.GetAllItems()
	for _,v := range items{
		fmt.Println(v.ItemId.Hex())
		fmt.Println(v.Title)
		fmt.Println(v.Description)
	}
}
