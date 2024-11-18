package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

// NotificationPayload represents the data sent to Gotify
type Notification struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	Extras   Extras `json:"extras,omitempty"`
	Sent     bool   `json:"sent"`
}

// Extras defines additional metadata for Gotify notifications
type Extras struct {
	Sound    string `json:"client::notification::sound,omitempty"`
	Icon     string `json:"client::notification::bigImage,omitempty"`
	Subtitle string `json:"client::notification::subtitle,omitempty"`
}

// In-memory queue to store unsent notifications
type NotificationQueue struct {
	mu            sync.Mutex
	notifications []Notification
}

// Add a notification to the queue
func (q *NotificationQueue) AddNotification(n Notification) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.notifications = append(q.notifications, n)
}

// Get all unsent notifications
func (q *NotificationQueue) getUnsentNotifications() []Notification {
	q.mu.Lock()
	defer q.mu.Unlock()
	var unsent []Notification
	for _, n := range q.notifications {
		if !n.Sent {
			unsent = append(unsent, n)
		}
	}
	return unsent
}

// Mark a notification as sent
func (q *NotificationQueue) markAsSent(id uint) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, n := range q.notifications {
		if n.ID == id {
			q.notifications[i].Sent = true
			break
		}
	}
}

func sendNotificationToDevices(notification Notification) error {
	gotifyURL := os.Getenv("GOTIFY_SERVER") // Your Gotify server URL
	appToken := os.Getenv("GOTIFY_TOKEN")   // Application Token from Gotify

	// Serialize payload to JSON
	jsonPayload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	// Send notification to Gotify server
	req, err := http.NewRequest("POST", gotifyURL+"/message", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", appToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log response for debugging
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification: %s, Response: %s", resp.Status, body)
	}

	log.Infof("Notification sent to all devices!")
	return nil
}

// Worker to process the queue
func (q *NotificationQueue) ProcessQueue() {
	for {
		unsentNotifications := q.getUnsentNotifications()

		if len(unsentNotifications) > 0 {
			for _, n := range unsentNotifications {
				err := sendNotificationToDevices(n)
				if err != nil {
					// Wait a few seconds before retrying
					log.Errorf("Failed to send notification: %v", err)
					time.Sleep(5 * time.Second)
				} else {
					q.markAsSent(n.ID)
				}
			}
		} else {
			log.Debug("No unsent notifications in the queue.")
		}

		// Sleep for a short time before checking again
		time.Sleep(10 * time.Second)
	}
}

var GotifyQueue = &NotificationQueue{}
