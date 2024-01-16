package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client, anl *Analytics) error

const (
	EventVisit = "visit"
	EventView = "view"
	EventUpdateAdmin = "uadmin"
)

type SendEventVisit struct {
	Sauce string `json:"sauce"`
	Agent string `json:"agent"`
}

type SendAdminUpdate struct {
	Visits  map[string]*Visit
}

func SendVisitHandler(event Event, client *Client, anl *Analytics) error{
	var payload SendEventVisit
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	analizer.addVisit(Visit{ Id: client.id, Ip: client.socket.RemoteAddr().String(), Views: 0, Duration: 0, Sauce: payload.Sauce, Agent: payload.Agent, Date: time.Now() })

	update := SendAdminUpdate{
		Visits: analizer.visits,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventUpdateAdmin

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room == "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil

}


func SendViewHandler(event Event, client *Client, anl *Analytics) error{
	analizer.updateViews(client.id)

	update := SendAdminUpdate{
		Visits: analizer.visits,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventUpdateAdmin

	for client := range client.manager.clients {
		// Only send to clients inside the same chatroom
		if client.room == "admin" {
			client.egress <- outgoingEvent
		}

	}
	return nil
}
