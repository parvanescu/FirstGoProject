package useCase

import (
	"ExGabi/interfaces"
	"ExGabi/model"
	"fmt"
)

//type IUseCase interfaces {
//	AddItem(id int,title string,description string)
//	DeleteItem(id int)model.IToDoItem
//	UpdateItem(id int,newTitle string,newDescription string)model.IToDoItem
//	GetItem(id int)model.IToDoItem
//	GetAll()//[]model.ToDoItem
//}

type UseCase struct{
	repository interfaces.IRepository
}

func New(repo interfaces.IRepository)interfaces.IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)Init(repository interfaces.IRepository){
	uC.repository = repository
}

func (uC *UseCase)AddItem(id int,title string,description string){
	var item interfaces.IToDoItem = model.New(id,title,description)
	uC.repository.Add(item)
}

func (uC *UseCase)DeleteItem(id int)interfaces.IToDoItem{
	return uC.repository.Delete(id)
}

func (uC *UseCase)UpdateItem(id int,newTitle string,newDescription string)interfaces.IToDoItem{
	return uC.repository.Update(id,model.New(id,newTitle,newDescription))
}

func(uC *UseCase)GetItem(id int)interfaces.IToDoItem{
	
	for _,v := range uC.repository.GetAll(){
		if v.GetItemId() == id{
			return v
		}
	}
	return model.New(0,"","")
}

func(uC *UseCase)GetAll(){
	items := uC.repository.GetAll()
	for _,v := range items{
		fmt.Println(v)
	}
}
