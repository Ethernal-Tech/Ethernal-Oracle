package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"oracle-test/plugins"
	"os"

	"github.com/joho/godotenv"
)

type LiveScore struct {
	api_key string
	address string
}

func (s *LiveScore) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s.address = os.Getenv("LIVESCORE_ADDRESS")
	s.api_key = os.Getenv("LIVESCORE_API_KEY")
}

func (s *LiveScore) GetMethods() []plugins.Method {
	return plugins.DefaulGetMethods(s)
}

func (s *LiveScore) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(s, methodName, params...)
}

func (s *LiveScore) Get_All_Leagues() ([]byte, error) {
	apiUrl := s.address + "list?Category=soccer"

	request, err := http.NewRequest("GET", apiUrl, nil)
	request.Header.Add("X-RapidAPI-Key", s.api_key)
	request.Header.Add("X-RapidAPI-Host", "livescore6.p.rapidapi.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}

	return body, nil
}

func main() {}

var ExportPlugin = LiveScore{}
