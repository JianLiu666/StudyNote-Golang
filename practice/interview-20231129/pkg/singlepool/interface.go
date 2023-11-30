package singlepool

import (
	"interview20231129/model"
	"interview20231129/pkg/e"
)

type SinglePool interface {
	AddSinglePersonAndMatch(user *model.User) e.CODE
	RemoveSinglePerson(name string) e.CODE
}
