package graph

import "ExGabi/useCase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	UseCase useCase.IUseCase
}
