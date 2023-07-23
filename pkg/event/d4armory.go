package event

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
)

type d4armoryData struct {
	Boss     Boss       `json:"boss,omitempty"`
	Helltide Helltide   `json:"helltide,omitempty"`
	Legion   Legion     `json:"legion,omitempty"`
	Whispers []Whispers `json:"whispers,omitempty"`
}
type Boss struct {
	Name             string `json:"name,omitempty"`
	ExpectedName     string `json:"expectedName,omitempty"`
	NextExpectedName string `json:"nextExpectedName,omitempty"`
	Timestamp        int    `json:"timestamp,omitempty"`
	Expected         int    `json:"expected,omitempty"`
	NextExpected     int    `json:"nextExpected,omitempty"`
	Territory        string `json:"territory,omitempty"`
	Zone             string `json:"zone,omitempty"`
}
type Helltide struct {
	Timestamp int    `json:"timestamp,omitempty"`
	Zone      string `json:"zone,omitempty"`
	Refresh   int    `json:"refresh,omitempty"`
}
type Legion struct {
	Timestamp    int    `json:"timestamp,omitempty"`
	Territory    string `json:"territory,omitempty"`
	Zone         string `json:"zone,omitempty"`
	Expected     int    `json:"expected,omitempty"`
	NextExpected int    `json:"nextExpected,omitempty"`
}
type Whispers struct {
	Quest int `json:"quest,omitempty"`
	End   int `json:"end,omitempty"`
}

func getD4EventData() (*d4armoryData, error) {
	log.Debugf("call d4armory.io API")
	// Get request
	resp, err := http.Get("https://d4armory.io/api/events/recent")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	var result *d4armoryData
	if err = json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		return nil, err
	}

	return result, err
}
