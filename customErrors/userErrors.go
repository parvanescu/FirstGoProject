package customErrors

type UserNotFoundError struct {
	Msg string
}

func (*UserNotFoundError) Error() string{
	return "User not found"
}

func NewUserNotFoundError() *UserNotFoundError{
	return &UserNotFoundError{"User not found."}
}
