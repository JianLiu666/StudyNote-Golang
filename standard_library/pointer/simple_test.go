package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectOperation(t *testing.T) {
	ast := assert.New(t)
	obj := CreateObject()

	ast.Equal("", obj.UUID)

	obj.SetUUIDByPointer("123")
	ast.Equal("123", obj.UUID)

	obj.SetUUIDByValue("456")
	ast.Equal("123", obj.UUID)

	obj = obj.SetUUIDByValue("456")
	ast.Equal("456", obj.UUID)
}

func TestNetestObjectOperation(t *testing.T) {
	ast := assert.New(t)
	nestedObj := CreateNestedObject()

	ast.Equal("", nestedObj.obj.UUID)

	nestedObj.obj.SetUUIDByPointer("123")
	ast.Equal("123", nestedObj.obj.UUID)
	ast.Equal("123", GetNestedObjectUUID(nestedObj))

	SetNestedObjectUuidByValue(nestedObj, "456")
	ast.Equal("123", nestedObj.obj.UUID)

	SetNestedObjectUuidByPointer(&nestedObj, "456")
	ast.Equal("456", nestedObj.obj.UUID)

	vObj := nestedObj.GetObjectByValue()
	vObj.SetUUIDByPointer("789")
	ast.Equal("456", nestedObj.obj.UUID)

	pObj := nestedObj.GetObjectByPointer()
	pObj.SetUUIDByPointer("789")
	ast.Equal("789", nestedObj.obj.UUID)
}
