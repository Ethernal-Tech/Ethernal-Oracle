package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oracle-test/plugins"
	"os"

	"github.com/joho/godotenv"
)

type Match struct {
	AwayTeam       string  `json:"away_team"`
	Bookie         string  `json:"bookie"`
	Competition    string  `json:"competition"`
	Country        string  `json:"country"`
	Date           string  `json:"date"`
	HomeTeam       string  `json:"home_team"`
	Match          string  `json:"match"`
	MatchStatus    string  `json:"match_status"`
	MatchTimestamp float64 `json:"match_timestamp"`
	MatchURL       string  `json:"match_url"`
	MatchID        string  `json:"matchid"`
	Sport          string  `json:"sport"`
	Time           string  `json:"time"`
}

type Quotas struct {
	Away           float64     `json:"away"`
	AwayTeam       interface{} `json:"away_team"`
	BScoreN        interface{} `json:"b_score_n"`
	BScoreY        interface{} `json:"b_score_y"`
	BetradarID     interface{} `json:"betradar_id"`
	Bookie         interface{} `json:"bookie"`
	Competition    interface{} `json:"competition"`
	Country        interface{} `json:"country"`
	Date           interface{} `json:"date"`
	DoubleChance12 interface{} `json:"double_chance_12"`
	DoubleChance1X interface{} `json:"double_chance_1X"`
	DoubleChance2X interface{} `json:"double_chance_2X"`
	Draw           interface{} `json:"draw"`
	DrawNoBet1     interface{} `json:"draw_no_bet_1"`
	DrawNoBet2     interface{} `json:"draw_no_bet_2"`
	FirstG1        interface{} `json:"first_g_1"`
	FirstG2        interface{} `json:"first_g_2"`
	FirstGX        interface{} `json:"first_g_X"`
	FirstH1        interface{} `json:"first_h_1"`
	FirstH2        interface{} `json:"first_h_2"`
	FirstHX        interface{} `json:"first_h_X"`
	Hand01_1       interface{} `json:"hand01_1"`
	Hand01_2       interface{} `json:"hand01_2"`
	Hand01X        interface{} `json:"hand01_X"`
	Hand02_1       interface{} `json:"hand02_1"`
	Hand02_2       interface{} `json:"hand02_2"`
	Hand02X        interface{} `json:"hand02_X"`
	Hand03_1       interface{} `json:"hand03_1"`
	Hand03_2       interface{} `json:"hand03_2"`
	Hand03X        interface{} `json:"hand03_X"`
	Hand10_1       interface{} `json:"hand10_1"`
	Hand10_2       interface{} `json:"hand10_2"`
	Hand10X        interface{} `json:"hand10_X"`
	Hand20_1       interface{} `json:"hand20_1"`
	Hand20_2       interface{} `json:"hand20_2"`
	Hand20X        interface{} `json:"hand20_X"`
	Hand30_1       interface{} `json:"hand30_1"`
	Hand30_2       interface{} `json:"hand30_2"`
	Hand30X        interface{} `json:"hand30_X"`
	Home           interface{} `json:"home"`
	HomeTeam       interface{} `json:"home_team"`
	LastG1         interface{} `json:"last_g_1"`
	LastG2         interface{} `json:"last_g_2"`
	LastGX         interface{} `json:"last_g_X"`
	Match          interface{} `json:"match"`
	MatchStatus    interface{} `json:"match_status"`
	MatchTimestamp interface{} `json:"match_timestamp"`
	MatchURL       interface{} `json:"match_url"`
	MatchID        interface{} `json:"matchid"`
	ScrapedDate    interface{} `json:"scraped_date"`
	Sport          interface{} `json:"sport"`
	Time           interface{} `json:"time"`
	TotalGoalsEven interface{} `json:"total_goals_even"`
	TotalGoalsOdd  interface{} `json:"total_goals_odd"`
	TotalOver05    interface{} `json:"total_over_05"`
	TotalOver15    interface{} `json:"total_over_15"`
	TotalOver25    interface{} `json:"total_over_25"`
	TotalOver35    interface{} `json:"total_over_35"`
	TotalOver45    interface{} `json:"total_over_45"`
	TotalUnder05   interface{} `json:"total_under_05"`
	TotalUnder15   interface{} `json:"total_under_15"`
	TotalUnder25   interface{} `json:"total_under_25"`
	TotalUnder35   interface{} `json:"total_under_35"`
	TotalUnder45   interface{} `json:"total_under_45"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type LiveScore struct {
	api_key string
	address string
}

func (s *LiveScore) Initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file %v", err)
	}

	s.address = os.Getenv("BET365_ADDRESS")
	s.api_key = os.Getenv("BET365_API_KEY")

	return nil
}

func (s *LiveScore) GetMethods() ([]plugins.Method, error) {
	return plugins.DefaultGetMethods(s)
}

func (s *LiveScore) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	return plugins.DefaultCallMethod(s, methodName, params...)
}

func (s *LiveScore) Get_Premier_League_Quotas() ([]byte, error) {
	matchesUrl := s.address + "matches_bet365?sport=soccer&country=england&competition=premier-league&match_urls=false"

	request, err := http.NewRequest("GET", matchesUrl, nil)
	request.Header.Add("X-RapidAPI-Key", s.api_key)
	request.Header.Add("X-RapidAPI-Host", "bet36528.p.rapidapi.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err)
	}

	var matches map[int32]Match
	err = json.Unmarshal(body, &matches)
	if err != nil {
		var errorMsg ErrorMessage
		err = json.Unmarshal(body, &errorMsg)
		if err == nil {
			return nil, fmt.Errorf(errorMsg.Message)
		}
		return nil, fmt.Errorf("Error decoding JSON: %s", err)
	}

	var quotas = make(map[string]Quotas, len(matches))

	for _, m := range matches {
		quotasUrl := s.address + fmt.Sprintf("odds_bet365?matchid=%s", m.MatchID)

		request, _ := http.NewRequest("GET", quotasUrl, nil)
		request.Header.Add("X-RapidAPI-Key", s.api_key)
		request.Header.Add("X-RapidAPI-Host", "bet36528.p.rapidapi.com")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("Error making request: %s", err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("Error making request: %s", err)
		}

		var quota map[string]Quotas
		err = json.Unmarshal(body, &quota)
		if err != nil {
			return nil, fmt.Errorf("Error decoding JSON: %s", err)
		}
		quotas[m.MatchID] = quota["0"]

		// limit to one request because of API daily limit
		break
	}

	return body, nil
}

func main() {}

var ExportPlugin = LiveScore{}
