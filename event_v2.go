package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Event includes the incident/alert details
type V2Event struct {
	RoutingKey string         `json:"routing_key"`
	Action     string         `json:"event_action"`
	DedupKey   string         `json:"dedup_key,omitempty"`
	Images     []V2EventImage `json:"images,omitempty"`
	Links      []V2EventLink  `json:"links,omitempty"`
	Client     string         `json:"client,omitempty"`
	ClientURL  string         `json:"client_url,omitempty"`
	Payload    *V2Payload     `json:"payload,omitempty"`
}

// EventLink is a link to be included on an event
type V2EventLink struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// EventImage is an image to be included, with an optional link and alt text
type V2EventImage struct {
	Src  string `json:"src"`
	Href string `json:"href,omitempty"`
	Alt  string `json:"alt,omitempty"`
}

// Payload represents the individual event details for an event
type V2Payload struct {
	Summary   string      `json:"summary"`
	Source    string      `json:"source"`
	Severity  string      `json:"severity"`
	Timestamp string      `json:"timestamp,omitempty"`
	Component string      `json:"component,omitempty"`
	Group     string      `json:"group,omitempty"`
	Class     string      `json:"class,omitempty"`
	Details   interface{} `json:"custom_details,omitempty"`
}

// Response is the json response body for an event
type V2EventResponse struct {
	RoutingKey  string `json:"routing_key"`
	DedupKey    string `json:"dedup_key"`
	EventAction string `json:"event_action"`
}

const v2eventEndPoint = "https://events.pagerduty.com/v2/enqueue"

// ManageEvent handles the trigger, acknowledge, and resolve methods for an event
func ManageEvent(e V2Event) (*V2EventResponse, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", v2eventEndPoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("HTTP Status Code: %d, Message: %s", resp.StatusCode, string(bytes))
	}
	var eventResponse V2EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}
	return &eventResponse, nil
}
