package utils

import (
	"ExGabi/payload"
	"errors"
	"regexp"
	"strings"
)

func CheckRegisterCredentials(user * payload.User)error{
	if len(user.FirstName) < 3{
		return errors.New("first name's length is too small")
	}
	if len(user.LastName) < 3{
		return errors.New("last name's length is to small")
	}
	if strings.Contains(user.Email,"@") == false{
		return errors.New("invalid email address")
	}
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+=?^_`{|}~-]+@[a-zA-Z0-9.]*.[a-z]+$")
	if !emailRegex.MatchString(user.Email){
		return errors.New("invalid email address")
	}
	return nil
}

func CheckPassword(password string)error{
	if len(password) < 6{
		return errors.New("password's length is to small")
	}
	return nil
}






