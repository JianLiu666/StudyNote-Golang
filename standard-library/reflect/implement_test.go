package main

import (
	"fmt"
	"reflect"
	"testing"
)

type CustomError struct {
}

func (*CustomError) Error() string {
	return ""
}

func Test_Implement(t *testing.T) {
	typeOfError := reflect.TypeOf((*error)(nil)).Elem()
	customErrorPtr := reflect.TypeOf(&CustomError{})
	customError := reflect.TypeOf(CustomError{})

	fmt.Println(customErrorPtr.Implements(typeOfError))
	fmt.Println(customError.Implements(typeOfError))
}
