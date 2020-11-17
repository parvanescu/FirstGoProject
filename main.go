package main

import (
	"ExGabi/interfaces"
	"ExGabi/repository"
	"ExGabi/useCase"
)

func main(){
	var repo interfaces.IRepository = repository.New()
	var uC interfaces.IUseCase = useCase.New(repo)
	uC.AddItem(1,"Title1","Description1")
	uC.AddItem(2,"Title2","Description2")
	uC.AddItem(3,"Title3","Description3")
	uC.AddItem(4,"Title4","Description4")
	uC.AddItem(5,"Title5","Description5")
	uC.GetAll()
}
