package repository

import (
	"ExGabi/interfaces"
	"ExGabi/model"
	"fmt"
)

//type IRepository interface{
//	Add(item model.IToDoItem)
//	Delete(id int) model.IToDoItem
//	Update(id int,newItem model.IToDoItem) model.IToDoItem
//	GetAll() []model.IToDoItem
//}

type Repository struct{
	dataArray []interfaces.IToDoItem
}

func New()interfaces.IRepository{
	fmt.Println("Repo initialized")
	repo :=Repository{make([]interfaces.IToDoItem,0,1)}
	return &repo
}


func (r *Repository)Add(item interfaces.IToDoItem){
	r.dataArray = append(r.dataArray,item)
}

func (r *Repository)Delete(id int)interfaces.IToDoItem{
	for i,v := range r.dataArray{
		if v.GetItemId() == id{
			item :=v
			r.dataArray = append(r.dataArray[0:i],r.dataArray[i+1:cap(r.dataArray)]...)
			return item
		}
	}
	return model.New(0,"","")
}

func (r *Repository)Update(id int,newItem interfaces.IToDoItem)interfaces.IToDoItem{
	for i,v  := range r.dataArray{
		if v.GetItemId() == id{
			newArray := append(r.dataArray[0:i],newItem)
			item := r.dataArray[i]
			r.dataArray = append(newArray,r.dataArray[i+1:cap(r.dataArray)]...)
			return item
			// Ce metoda ?????
			//r.dataArray[i].Description = newItem.Description
			//r.dataArray[i].Title = newItem.Title
			//r.dataArray[i].ItemId=newItem.ItemId
		}
	}
	return model.New(0,"","")
}

func (r *Repository)GetAll() []interfaces.IToDoItem{
	return r.dataArray
}
