package helpers

var (
	Constant = ""
	Variable = ""
)

func FuncFactory() func() {
	return func() {}
}

func StructFactory() struct{ Field string } {
	return struct{ Field string }{}
}

var Struct = struct{ StructFactory string }{}

type StructType struct{}

type InterfaceType interface{}
