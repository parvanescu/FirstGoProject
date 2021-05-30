package useCase

import (
	"ExGabi/customErrors"
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"ExGabi/utils"
	"ExGabi/utils/token"
	"errors"
)

type UseCase struct {
	mutationRepository mutations.IRepository
}



func New(repo mutations.IRepository) mutations.IUseCase {
	return &UseCase{repo}
}

func (uC *UseCase) AddItem(item *payload.Item) (*response.Item, error) {
	tokenClaims, err := token.CheckToken(item.Token)
	if err != nil {
		return nil, err
	}
	newToken, err := token.CreateToken(tokenClaims)
	if err != nil {
		return nil, nil
	}
	id, err := uC.mutationRepository.AddItem(tokenClaims.Id, item)
	if err != nil{
		return nil, err
	}
	return &response.Item{Id: id, Token: newToken}, nil
}
func (uC *UseCase) DeleteItem(item *payload.Item) (string, error) {
	tokenClaims, err := token.CheckToken(item.Token)
	if err != nil {
		return "", err
	}
	newToken, err := token.CreateToken(tokenClaims)
	if err != nil {
		return "", err
	}
	err = uC.mutationRepository.DeleteItem(tokenClaims.Id, item)
	return newToken, err
}
func (uC *UseCase) UpdateItem(item *payload.Item) (*response.Item, error) {
	tokenClaims, err := token.CheckToken(item.Token)
	if err != nil {
		return nil, err
	}
	newToken, err := token.CreateToken(tokenClaims)
	if err != nil {
		return nil, err
	}
	newItem, err := uC.mutationRepository.UpdateItem(tokenClaims.Id, item)
	if err != nil {
		return nil, err
	}
	newItem.Token = newToken
	return newItem, nil
}

func (uC *UseCase) DeleteUser(user *payload.User) (string, error) {
	tkn, err := token.CheckToken(user.Token)
	if err != nil {
		return "", err
	}
	newToken, err := token.CreateToken(tkn)
	if err != nil {
		return "", err
	}
	err = uC.mutationRepository.DeleteUser(user)
	return newToken, err
}
func (uC *UseCase) UpdateUser(user *payload.User) (*response.User, error) {
	tknClaims, err := token.CheckToken(user.Token)
	if err != nil {
		return nil, err
	}
	newToken, err := token.CreateToken(tknClaims)
	if err != nil {
		return nil, err
	}
	user.Id = tknClaims.Id
	newUser, err := uC.mutationRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	newUser.Token = newToken
	return newUser, nil
}
func (uC *UseCase) UpdateUserPerformedByLeader(user *payload.User,oldEmail string)(*response.User,error){
	tknClaims, err := token.CheckToken(user.Token)
	if err != nil {
		return nil, err
	}
	newToken, err := token.CreateToken(tknClaims)
	if err != nil {
		return nil, err
	}
	oldUser,err := uC.mutationRepository.GetUserByEmail(&payload.User{Email: oldEmail})
	if err != nil{
		return nil, err
	}

	user.Id = oldUser.Id
	user.OrganisationId = oldUser.OrganisationId
	user.Status = oldUser.Status
	if user.LastName == ""{user.LastName = oldUser.LastName}
	if user.FirstName == ""{user.FirstName = oldUser.FirstName}
	if user.Email == ""{user.Email = oldEmail}
	if user.Password == ""{user.Password = oldUser.Password}

	responseUser,err := uC.mutationRepository.UpdateUser(user)
	if err!=nil{
		return nil,err
	}
	responseUser.Token = newToken
	return responseUser, nil

}

func (uC *UseCase) Register(user *payload.User, organisation *payload.Organisation) (*response.User, *response.Organisation, error) {
	err := utils.CheckRegisterCredentials(user)
	if err != nil {
		return nil, nil, err
	}

	dbUser, err := uC.mutationRepository.GetUserByEmail(user)
	if err != nil {
		if err.Error() == "not found" {
			dbOrganisation, err := uC.mutationRepository.GetOrganisationByCUI(organisation)
			switch err.(type) {
			case *customErrors.OrganisationNotFoundError:
				userId, organisationId, err := uC.mutationRepository.AddUserAndOrganisation(user, organisation)
				if err != nil {
					return nil, nil, err
				}
				return &response.User{Id: userId}, &response.Organisation{Id: organisationId}, nil
			case nil:
				if dbOrganisation.Status == true {
					return nil, nil, errors.New("please talk to the organisation owner to add your credentials from the application and then login")
				} else {
					return nil, nil, errors.New("the organisation already exist and it is bounded to another email please talk to the site owners")
				}
			default:
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	dbOrganisation, err := uC.mutationRepository.GetOrganisationByCUI(organisation)
	switch err.(type) {
	case *customErrors.OrganisationNotFoundError:
		return nil, nil, errors.New("this account already exists")
	case nil:
		if dbUser.Status == false && dbOrganisation.Status == false {
			return &response.User{Id: dbUser.Id}, &response.Organisation{Id: dbOrganisation.Id}, nil
		} else {
			return nil, nil, errors.New("the account or the organisation are already active and used please check the credentials")
		}
	default:
		return nil, nil, err
	}

}
func (uC *UseCase) Login(user *payload.User) (*response.User, error) {
	responseUser, err := uC.mutationRepository.GetUserByCredentials(user)
	if err != nil {
		return nil, err
	}

	if responseUser.Status == false{
		return &response.User{Id: responseUser.Id,OrganisationId: responseUser.OrganisationId,Status: false}, nil
	}

	tkn, err := token.CreateToken(&response.User{Id: responseUser.Id, Email: responseUser.Email, OrganisationId: responseUser.OrganisationId})
	if err != nil {
		return nil, err
	}
	return &response.User{Token: tkn,Status: true}, nil
}
func (uC *UseCase) AddInactiveUser(user *payload.User)(*response.User,error){
	tknClaims, err := token.CheckToken(user.Token)
	if err != nil {
		return nil, err
	}
	newToken, err := token.CreateToken(tknClaims)
	if err != nil {
		return nil, err
	}
	user.Id = tknClaims.Id

	err = utils.CheckRegisterCredentials(user)
	if err != nil {
		return nil, err
	}

	_, err = uC.mutationRepository.GetUserByEmail(user)
	if err == nil{
		return nil, errors.New("email already exists")
	}
	if err.Error() == "not found"{
		user.Status=false
		_,err:=uC.mutationRepository.AddUser(tknClaims.OrganisationId,user)
		if err!=nil{
			return nil, err
		}
		return &response.User{Token: newToken},nil
	} else {return nil, errors.New("DB error")}

}

func (uC *UseCase) SetUserPassword(user *payload.User) error {
	err := utils.CheckPassword(user.Password)
	if err != nil {
		return err
	}
	err = uC.mutationRepository.UpdateUserPassword(user, &payload.Organisation{Id: user.OrganisationId})
	if err != nil {
		return err
	}
	return nil
}

func (uC *UseCase) GetMatchingSearch(item *payload.Item) (*[]response.Item, string, error) {
	tokenClaims, err := token.CheckToken(item.Token)
	if err != nil {
		return nil, item.Token, err
	}
	newToken, err := token.CreateToken(tokenClaims)
	if err != nil {
		return nil, item.Token, err
	}
	descriptionSearchCriteria := &payload.Item{Description: item.Title}

	descriptionItems, err := uC.mutationRepository.GetMatchingItems(tokenClaims.Id, descriptionSearchCriteria)

	if err != nil {
		emptyList := new([]response.Item)
		return emptyList, newToken, err
	}

	return descriptionItems, newToken, nil
}


func (uC *UseCase) AddPosition(position *payload.Position) (*response.Position, error) {
	tokenClaims, err := token.CheckToken(position.Token)
	if err != nil{
		return nil, err
	}
	newToken ,err := token.CreateToken(tokenClaims)
	position.OrganisationId = tokenClaims.OrganisationId
	id,err := uC.mutationRepository.AddPositionToOrganisation(position)
	if err != nil{
		return nil, err
	}
	return &response.Position{Id: id,Token: newToken}, nil
}

func (uC *UseCase) ExchangePositionsAccessLevel(bottomPosition *payload.Position, topPosition *payload.Position) (*response.Position, *response.Position, error) {
	tokenClaims, err := token.CheckToken(bottomPosition.Token)
	if err != nil{
		return nil,nil, err
	}
	newToken ,err := token.CreateToken(tokenClaims)
	newTopPosition,newBottomPosition,err := uC.mutationRepository.ExchangePositions(bottomPosition,topPosition)
	if err != nil{
		return nil, nil, err
	}
	newTopPosition.Token=newToken
	newBottomPosition.Token = newToken
	return newBottomPosition, newTopPosition, nil
}
//func createIndexes(collectionName string, coll *mongo.Collection) error {
//	if collectionName == "quote_requests" {
//		mod1 := mongo.IndexModel{Keys: bson.M{"origin_name": 1}, Options: nil}
//		mod2 := mongo.IndexModel{Keys: bson.M{"destination_name": 1}, Options: nil}
//		_, err := coll.Indexes().CreateOne(context.TODO(), mod1)
//		if err != nil {
//			fmt.Println("fail create index for", collectionName)
//			return err
//		}
//		_, err = coll.Indexes().CreateOne(context.TODO(), mod2)
//		if err != nil {
//			fmt.Println("fail create index for", collectionName)
//			return err
//		}
//	}
//	return nil
//}
