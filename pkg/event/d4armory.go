package event

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

// Host API
var Host string

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

func InitHost() {
	// Init host value
	Host = os.Getenv("D4ARMORY_HOST")
	mockHost := os.Getenv("MOCK_D4ARMORY_HOST")
	if len(mockHost) > 0 {
		Host = mockHost
	}
	log.Debugf("d4armory.io host %s", Host)
}

func getD4EventData() (*d4armoryData, error) {
	log.Debugf("call d4armory.io API")
	// build url
	url := fmt.Sprintf("%s/api/events/recent", Host)
	// Get request
	resp, err := http.Get(url)
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
