package singlepool

import (
	"interview20231129/model"
	"interview20231129/pkg/e"
)

type QueryOpts struct {
	Name        string
	MinHeight   int
	MaxHeight   int
	Gender      int
	MinNumDates int
	MaxNumDates int
}

type SinglePool interface {
	AddSinglePersonAndMatch(user *model.User) e.CODE
	RemoveSinglePerson(name string) e.CODE
	QuerySinglePeople(limit int, opts *QueryOpts) ([]*model.User, e.CODE)
}
