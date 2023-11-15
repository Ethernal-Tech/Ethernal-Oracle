package plugins

type IPlugin interface {
	Initialize()
	CallMethod(methodName string) string
}
