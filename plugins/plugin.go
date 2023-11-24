package plugins

type Method struct {
	MethodName   string
	InputParams  []Param
	OutputParams []Param
}

type Param struct {
	ParamName string
	ParamType string
}

type IPlugin interface {
	Initialize()
	GetMethods() []Method
	CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error)
}
