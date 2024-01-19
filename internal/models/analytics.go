package models

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 4096
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

func checkOrigin(r *http.Request) bool {

	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	// Update this to HTTPS
	case os.Getenv("HOST"):
		return true
	default:
		return false
	}
}

var upgrader = websocket.Upgrader{ CheckOrigin: checkOrigin, ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type Analytics struct {
	visits map[string]*Visit
	lock sync.Mutex
}

var analizer = Analytics { visits: make(map[string]*Visit) };

func (anl *Analytics) addVisit(visit Visit) {
	anl.lock.Lock()
	defer anl.lock.Unlock()
	anl.visits[visit.Id] = &visit
}

func (anl * Analytics) updateViews(id string) {
	anl.lock.Lock()
	defer anl.lock.Unlock()
	anl.visits[id].Views += 1
}

func (anl * Analytics) archiveVisit(id string) {
	anl.lock.Lock()
	defer anl.lock.Unlock()

	visit := anl.visits[id]
	if visit == nil {
		log.Debug("Visit not found")
		return
	}
	defer delete(anl.visits, id)

	if err := visit.Archive(); err != nil {
		log.Error(err)
		return
	}
}

type ConnectionManager struct {
	clients map[*Client]bool
	connect chan *Client
	disconnect chan *Client
	handlers map[string]EventHandler
	// otps is a map of allowed OTP to accept connections from
	otps RetentionMap
}

func (cm *ConnectionManager) GenerateNewOtp() string {
	return cm.otps.NewOTP().Key
}

func NewManager(ctx context.Context) *ConnectionManager {
	cm :=  &ConnectionManager{
		connect: make(chan *Client),
		disconnect: make(chan *Client),
		clients: make(map[*Client]bool),
		handlers: make(map[string]EventHandler),
		otps: NewRetentionMap(ctx, 5*time.Second),
	}

	cm.setupEventHandlers()



	return cm
}

// setupEventHandlers configures and adds all handlers
func (m *ConnectionManager) setupEventHandlers() {
	m.handlers[EventVisit] = SendVisitHandler
	m.handlers[EventView] = SendViewHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *ConnectionManager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c, &analizer); err != nil {
			return err
		}
		return nil
	} else {
		return  errors.New("this event type is not supported")
	}
}

func (cm *ConnectionManager) Run() {
	for {
		select {
		case client := <- cm.connect:
			cm.clients[client] = true
		case client := <- cm.disconnect:
			if _, ok := cm.clients[client]; ok {
					close(client.egress)
					analizer.archiveVisit(client.id)
					client.socket.Close()
					delete(cm.clients, client)
			}
		}
	}
}

func (cm *ConnectionManager) ServeWS(c echo.Context) error {
	socket, err := upgrader.Upgrade(c.Response() , c.Request(), nil)
	if err != nil {
		log.Fatal("Serve HTTP Sockets Error: ", err)
		return err
	}

	log.Info("Connection Received")

	client := &Client{
		id: uuid.NewV4().String(),
		socket: socket,
		egress: make(chan Event, messageBufferSize),
		manager: cm,
		room: "base",
		sauce: c.Request().Header.Get("Referer"),
		agent: c.Request().Header.Get("User-Agent"),
	}

	cm.connect <- client

	go client.read()

	go client.write()

	return nil
}
type Client struct {
	id string

	socket *websocket.Conn

	egress chan Event

	manager *ConnectionManager

	room string

	sauce string

	agent string
}

func (client * Client) read() {
	defer func () {
		client.manager.disconnect <- client
	}()

	client.socket.SetReadLimit(messageBufferSize)

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := client.socket.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error(err)
		return
	}

	client.socket.SetPongHandler(client.pongHandler)

	for {
		_, payload, err := client.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("error reading message: %v", err)
			}
			return
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Errorf("error marshalling message: %v", err)
			return // Breaking the connection here might be harsh xD
		}

		log.Infof("event received: %v", request)

		if err := client.manager.routeEvent(request, client); err != nil {
			log.Errorf("Error handeling Message: ", err)
		}
	}
}

func (client * Client) write() {
	ticker := time.NewTicker(pingInterval)

	defer func () {
			ticker.Stop()
		 	client.manager.disconnect <- client
	}()

	for {
		select {
			case message, ok := <-client.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := client.socket.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Infof("connection closed: %v", err)
				}
				// Return to close the goroutine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Print(err)
				return // closes the connection, should we really
			}
			// Write a Regular text message to the connection
			if err := client.socket.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Print(err)
			}
			log.Print("sent message")
		case <-ticker.C:
			log.Print("ping")
			// Send the Ping
			if err := client.socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (client *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Print("pong")
	return client.socket.SetReadDeadline(time.Now().Add(pongWait))
}

