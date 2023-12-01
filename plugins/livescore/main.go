package main

import (
	"fmt"
	"io"
	"net/http"
	"oracle-test/plugins"
	"os"

	"github.com/joho/godotenv"
)

type LiveScore struct {
	api_key string
	address string
}

func (s *LiveScore) Initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file %v", err)
	}

	s.address = os.Getenv("LIVESCORE_ADDRESS")
	s.api_key = os.Getenv("LIVESCORE_API_KEY")

	return nil
}

func (s *LiveScore) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(s)
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
