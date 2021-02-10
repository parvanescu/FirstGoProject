package delivery

import (
	"ExGabi/payload"
	"ExGabi/useCase"
	"github.com/graphql-go/graphql"
)

type Handler struct {
	Uc useCase.IUseCase
}

func New(uc useCase.IUseCase) graphql.Fields {
	h := &Handler{Uc: uc}
	return graphql.Fields{
		"login": h.login(),
		
	}
}


func (h *Handler) login() *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Description: "Login",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{Type: graphql.String},
			"password": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			email := p.Args["email"].(string)
			password := p.Args["password"].(string)
			tkn,err:=h.Uc.Login(&payload.User{Email: email,Password: password})
			return tkn,err
		},
}
}