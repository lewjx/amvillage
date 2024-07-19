package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Authenticated bool
	TeamID        int

	ws     *websocket.Conn
	closed bool
	mu     *sync.Mutex
}

// handleWebsocket handles the upgrading of websocket before passing it onto the internal server.
func (state *GameState) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("error upgrading websocket: %w", err)
		return
	}
	conn := &Conn{
		Authenticated: false,
		TeamID:        -1,
		ws:            ws,
		mu:            new(sync.Mutex),
	}
	defer conn.Close()
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(maxPongGap))
	conn.ws.SetPongHandler(conn.pongHandler)
	for {
		_, message, err := conn.ws.ReadMessage()
		switch {
		case err == nil:
			// Process message below.
		case websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway):
			// Unexpected close.
			log.Printf("Unexpected close: %s\n", err)
			return
		default:
			return
		}
		cmd, err := ParseCommand(message)
		if err != nil {
			log.Printf("Error parsing command: %s\n", err)
			return
		}
		_, isLogin := cmd.(Auth)
		if isLogin != conn.Authenticated {
			log.Printf("Unexpected command %T: Authenticated is %v\n", cmd, conn.Authenticated, err)
			return
		}
		if auth, isLogin := cmd.(Auth); isLogin {
			for i, team := range state.Teams {
				if team.Secret == auth.Secret {
					conn.TeamID = i
					break
				}
			}
			if conn.TeamID >= 0 {
				continue
			}
			log.Printf("Invalid secret: %q\n", auth.Secret)
			return
		}
		// TODO: Authentication/Check permission
		state.msg <- cmd
	}
}

func (conn *Conn) Close() {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.closed = true
	conn.ws.Close()
}

func (conn *Conn) pingTicker() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for range ticker.C {
		conn.mu.Lock()
		if conn.closed {
			conn.mu.Unlock()
			return
		}
		conn.ws.WriteMessage(websocket.PingMessage, []byte{})
		conn.mu.Unlock()
	}
}

func (conn *Conn) pongHandler(string) error {
	conn.ws.SetReadDeadline(time.Now().Add(maxPongGap))
	return nil
}

const (
	maxMessageSize = 10000
	pingInterval   = time.Minute
	maxPongGap     = pingInterval + 5*time.Second
)

type (
	CommandType string
	CommandJSON struct {
		Type CommandType `json:"type"`
	}
)

const (
	CommandAuth               CommandType = "auth"
	CommandLock               CommandType = "lock"
	CommandTransfer           CommandType = "transfer"
	CommandNotice             CommandType = "notice"
	CommandNoticeStatusUpdate CommandType = "notice_status_update"
)

func ParseCommand(bytes []byte) (Command, error) {
	var cmd CommandJSON
	if err := json.Unmarshal(bytes, &cmd); err != nil {
		return nil, fmt.Errorf("error decoding message: %w", err)
	}
	var decoded Command
	switch cmd.Type {
	case CommandAuth:
		decoded = Auth{}
	case CommandLock:
		decoded = Lock{}
	case CommandTransfer:
		decoded = Transfer{}
	case CommandNotice:
		decoded = Notice{}
	case CommandNoticeStatusUpdate:
		decoded = NoticeStatusUpdate{}
	default:
		return nil, fmt.Errorf("unknown command type: %s", cmd.Type)
	}
	if err := json.Unmarshal(bytes, &cmd); err != nil {
		return nil, fmt.Errorf("error decoding message of type %q: %w", cmd.Type, err)
	}
	return decoded, nil
}

type GameStateJSON struct {
	Teams  []Team `json:"teams"`
	TeamID int    `json:"team_id"`
}

func (conn *Conn) Send(state *GameState) error {
	stateJSON := GameStateJSON{
		TeamID: conn.TeamID,
		Teams:  state.Teams,
	}
	b, err := json.Marshal(stateJSON)
	if err != nil {
		return fmt.Errorf("error sending message: cannot marshal message: %v", err)
	}
	return conn.ws.WriteMessage(websocket.TextMessage, b)
}
