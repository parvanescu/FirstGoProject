package useCase

import (
	"ExGabi/interfaces"
	"ExGabi/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UseCase struct{
	itemRepository interfaces.IRepository
	userRepository interfaces.IRepository
}

func New(itemRepo interfaces.IRepository,userRepo interfaces.IRepository)interfaces.IUseCase{
	return &UseCase{itemRepo,userRepo}
}

func (uC *UseCase)Init(repository interfaces.IRepository){
	uC.itemRepository = repository
}

func (uC *UseCase)AddItem(title string,description string){
	var item model.ToDoItem = model.New(title,description)
	uC.itemRepository.Add(item)
}

func (uC *UseCase)DeleteItem(id primitive.ObjectID)model.ToDoItem{
	return uC.itemRepository.Delete(id)
}

func (uC *UseCase)UpdateItem(id primitive.ObjectID,newTitle string,newDescription string)model.ToDoItem{
	return uC.itemRepository.Update(id,model.New(newTitle,newDescription))
}

func(uC *UseCase)GetItem(id primitive.ObjectID)model.ToDoItem{
	
	for _,v := range uC.itemRepository.GetAll(){
		if v.GetItemId() == id{
			return v
		}
	}
	return model.New("","")
}

func(uC *UseCase)GetAll(){
	items := uC.itemRepository.GetAll()
	for _,v := range items{
		fmt.Println(v.GetItemId().Hex())
		fmt.Println(v.GetTitle())
		fmt.Println(v.GetDescription())
	}
}
