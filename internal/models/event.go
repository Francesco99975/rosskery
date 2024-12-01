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

type HtmlData struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventVisit             = "visit"
	EventView              = "view"
	EventAuthAdmin         = "authadmin"
	EventUpdateVisitsAdmin = "uvadmin"
	EventSettingsChanged   = "settingschanged"
	EventNewProduct        = "newproduct"
	EventUpdateProduct     = "updateproduct"
	EventRemoveProduct     = "removeproduct"
	EventNewCategory       = "newcategory"
	EventRemoveCategory    = "removecategory"
	EventOrdersChanged     = "orderschanged"
	EventCustomersChanged  = "customerschanged"
)

func SendAdminUpdateHandler(event Event, client *Client) error {

	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = event.Type

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room == "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil
}

func SendVisitHandler(event Event, client *Client) error {
	var source string
	if len(client.sauce) > 0 {
		source = client.sauce
	} else {
		source = "direct"
	}

	analizer.addVisit(Visit{Id: client.id, Ip: client.ip, Views: 0, Duration: 0, Sauce: source, Agent: client.agent, Date: time.Now()})

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

func SendNewProductHandler(event Event, client *Client) error {
	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventNewProduct

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room != "admin" {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func SendUpdateProductHandler(event Event, client *Client) error {
	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventUpdateProduct

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room != "admin" {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func SendRemoveProductHandler(event Event, client *Client) error {
	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventRemoveProduct

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room != "admin" {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func SendNewCategoryHandler(event Event, client *Client) error {
	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventNewCategory

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room != "admin" {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func SendRemoveCategoryHandler(event Event, client *Client) error {
	var outgoingEvent Event
	outgoingEvent.Payload = event.Payload
	outgoingEvent.Type = EventRemoveCategory

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
