package singlepool

import "interview20231129/model"

type SinglePool interface {
	AddSinglePersonAndMatch(user *model.User)
	RemoveSinglePerson(name string)
}
