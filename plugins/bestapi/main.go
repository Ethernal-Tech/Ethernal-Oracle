package main

type BestApi struct{}

func (b *BestApi) Initialize() {

}

func (b *BestApi) CallMethod(methodName string) string {
	return "Hello from BestAPI"
}

func main() {}

var ExportPlugin = BestApi{}
