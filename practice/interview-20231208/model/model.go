package model

import "interview20231208/pkg/e"

type Order struct {
	RoleType     e.ROLE_TYPE     `json:"roleType"`     //
	OrderType    e.ORDER_TYPE    `json:"orderType"`    //
	DurationType e.DURATION_TYPE `json:"durationType"` //
	Price        int             `json:"price"`        //
	Quantity     int             `json:"quantity"`     //
	Status       int             `json:"status"`       //
	Timestamp    int             `json:"timestamp"`    //
}
