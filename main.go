package main

import (
	"ExGabi/repository"
	"ExGabi/useCase"
)

func main(){
	var repo repository.IRepository = repository.New()
	var uC useCase.IUseCase = useCase.New(repo)
	uC.AddItem(1,"Title1","Description1")
	uC.AddItem(2,"Title2","Description2")
	uC.AddItem(3,"Title3","Description3")
	uC.AddItem(4,"Title4","Description4")
	uC.AddItem(5,"Title5","Description5")
	uC.GetAll()
	//repo.Add(toDo)
	//toDo2 := model.ToDoItem{
	//	ItemId: 45,
	//	Title: "Titlu2",
	//	Description: "Descr2",
	//}
	//fmt.Println(repo.GetAll())
	//repo.Delete(30)
	//fmt.Println(repo.GetAll())
	//repo.Add(toDo)
	//repo.Update(30,toDo2)
	//fmt.Println(repo.GetAll())
}
