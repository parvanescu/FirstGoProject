package customErrors

type OrganisationNotFoundError struct{
	Msg string
}

func (org *OrganisationNotFoundError) Error() string{
	return org.Msg
}

func NewOrganisationNotFoundError() *OrganisationNotFoundError{
	return &OrganisationNotFoundError{"Organisation not found."}
}
