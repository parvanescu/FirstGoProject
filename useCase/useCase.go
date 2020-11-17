package useCase

import (
	"ExGabi/model"
	"ExGabi/repository"
	"fmt"
)

type IUseCase interface {
	AddItem(id int,title string,description string)
	DeleteItem(id int)model.IToDoItem
	UpdateItem(id int,newTitle string,newDescription string)model.IToDoItem
	GetItem(id int)model.IToDoItem
	GetAll()//[]model.ToDoItem
}

type UseCase struct{
	repository repository.IRepository
}

func New(repo repository.IRepository)IUseCase{
	return &UseCase{repo}
}

func (uC *UseCase)Init(repository repository.IRepository){
	uC.repository = repository
}

func (uC *UseCase)AddItem(id int,title string,description string){
	var item model.IToDoItem = model.New(id,title,description)
	uC.repository.Add(item)
}

func (uC *UseCase)DeleteItem(id int)model.IToDoItem{
	return uC.repository.Delete(id)
}

func (uC *UseCase)UpdateItem(id int,newTitle string,newDescription string)model.IToDoItem{
	return uC.repository.Update(id,model.New(id,newTitle,newDescription))
}

func(uC *UseCase)GetItem(id int)model.IToDoItem{
	
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
