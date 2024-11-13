package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventVisit             = "visit"
	EventView              = "view"
	EventAuthAdmin         = "authadmin"
	EventUpdateVisitsAdmin = "uvadmin"
	EventSettingsChanged   = "settingschanged"
)

func SendVisitHandler(event Event, client *Client) error {
	var source string
	if len(client.sauce) > 0 {
		source = client.sauce
	} else {
		source = "direct"
	}

	analizer.addVisit(Visit{Id: client.id, Ip: client.socket.RemoteAddr().String(), Views: 0, Duration: 0, Sauce: source, Agent: client.agent, Date: time.Now()})

	update := VisitsResponse{
		Current: len(analizer.visits),
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventUpdateVisitsAdmin

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room == "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil

}

func SendViewHandler(event Event, client *Client) error {
	analizer.updateViews(client.id)

	views, err := CountTotalViews()
	if err != nil {
		return fmt.Errorf("failed to get total views: %v", err)
	}

	update := VisitsResponse{
		TotalViews: views,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventUpdateVisitsAdmin

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room == "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil
}

func SendSettingsChangeHandler(event Event, client *Client) error {

	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventSettingsChanged

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room != "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil
}

type SendOtp struct {
	OTP string `json:"otp"`
}

func SendOtpHandler(event Event, client *Client) error {
	var payload SendOtp
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Verify OTP is existing
	if !client.manager.otps.VerifyOTP(payload.OTP) {
		return fmt.Errorf("authauthorized bad otp in request")
	}

	client.room = "admin"
	return nil
}
