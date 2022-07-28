package pointer

type Object struct {
	UUID        string
	Description string
	NumberValue int
}

func CreateObject() Object {
	return Object{}
}
func (o *Object) SetUUIDByPointer(uuid string) {
	o.UUID = uuid
}

func (o *Object) SetDescriptionByPointer(text string) {
	o.Description = text
}
func (o *Object) SetNumberValueByPointer(number int) {
	o.NumberValue = number
}

func (o Object) SetUUIDByValue(uuid string) Object {
	o.UUID = uuid
	return o
}

func (o Object) SetDescriptionByValue(text string) Object {
	o.Description = text
	return o
}
func (o Object) SetNumberValueByValue(number int) Object {
	o.NumberValue = number
	return o
}

type NestedObject struct {
	obj Object
}

func CreateNestedObject() NestedObject {
	return NestedObject{
		obj: CreateObject(),
	}
}

func (b *NestedObject) GetObjectByValue() Object {
	return b.obj
}

func (b *NestedObject) GetObjectByPointer() *Object {
	return &b.obj
}

func GetNestedObjectUUID(nestedObj NestedObject) string {
	return nestedObj.obj.UUID
}

func SetNestedObjectUuidByValue(nestedObj NestedObject, str string) {
	nestedObj.obj.SetUUIDByPointer(str)
}

func SetNestedObjectUuidByPointer(nestedObj *NestedObject, str string) {
	nestedObj.obj.SetUUIDByPointer(str)
}
