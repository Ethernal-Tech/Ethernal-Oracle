package plugins

type IPlugin interface {
	Initialize()
	CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error)
}
