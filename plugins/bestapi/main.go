package main

type BestApi struct{}

func (b *BestApi) Initialize() {

}

func (b *BestApi) CallMethod(methodName string, paramBytes ...[]byte) ([]byte, error) {
	return []byte("Hello from BestAPI"), nil
}

func main() {}

var ExportPlugin = BestApi{}
