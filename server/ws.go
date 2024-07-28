package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var currentUserID atomic.Int32

type Conn struct {
	Authenticated bool
	TeamID        int
	UserID        int
	Nickname      string

	ws     *websocket.Conn
	closed bool
	mu     *sync.Mutex
}

// handleWebsocket handles the upgrading of websocket before passing it onto the internal server.
func (state *GameState) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("error upgrading websocket: %v", err)
		return
	}
	userID := currentUserID.Add(1)
	conn := &Conn{
		Authenticated: false,
		TeamID:        -1,
		UserID:        int(userID),
		ws:            ws,
		mu:            new(sync.Mutex),
	}
	defer conn.Close()
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(maxPongGap))
	go conn.pingTicker()
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
		cmd, err := conn.ParseCommand(message)
		if err != nil {
			log.Printf("Error parsing command: %s\n", err)
			return
		}
		_, isLogin := cmd.(Auth)
		if isLogin == conn.Authenticated {
			log.Printf("Unexpected command %T: Authenticated is %v\n", cmd, conn.Authenticated)
			return
		}
		if auth, isLogin := cmd.(Auth); isLogin {
			for i, team := range state.Teams {
				if team.Secret == auth.Secret {
					conn.TeamID = i
					break
				}
			}
			if conn.TeamID < 0 {
				log.Printf("Invalid secret: %q\n", auth.Secret)
				return
			}
			if strings.TrimSpace(auth.Nickname) == "" {
				log.Printf("Invalid nickname\n")
				return
			}
			// Login successful.
			conn.Nickname = auth.Nickname
			conn.Authenticated = true
			conn.Send(state)
		}
		// TODO: Authentication/Check permission
		state.msg <- ConnCommand{
			Conn:    conn,
			Command: cmd,
		}
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
	CommandUnlock             CommandType = "unlock"
	CommandTransfer           CommandType = "transfer"
	CommandNotice             CommandType = "notice"
	CommandNoticeStatusUpdate CommandType = "notice_status_update"
)

func (conn *Conn) ParseCommand(bytes []byte) (Command, error) {
	var cmd CommandJSON
	if err := json.Unmarshal(bytes, &cmd); err != nil {
		return nil, fmt.Errorf("error decoding message: %w", err)
	}
	var decoded Command
	var err error
	switch cmd.Type {
	case CommandAuth:
		var auth Auth
		err = json.Unmarshal(bytes, &auth)
		decoded = auth
	case CommandLock:
		decoded = Lock{}
	case CommandUnlock:
		decoded = Unlock{}
	case CommandTransfer:
		var transfer Transfer
		err = json.Unmarshal(bytes, &transfer)
		decoded = transfer
	case CommandNotice:
		var notice Notice
		err = json.Unmarshal(bytes, &notice)
		decoded = notice
	case CommandNoticeStatusUpdate:
		var nsu NoticeStatusUpdate
		err = json.Unmarshal(bytes, &nsu)
		decoded = nsu
	default:
		return nil, fmt.Errorf("unknown command type: %s", cmd.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("error decoding message of type %q: %w", cmd.Type, err)
	}
	return decoded, nil
}

type GameStateJSON struct {
	// Type should always be "state".
	Type       string  `json:"type"`
	TeamID     int     `json:"team_id"`
	UserID     int     `json:"user_id"`
	Score      []Score `json:"score"`
	LockHolder []Queue `json:"lock_holder"`
}

type Score struct {
	Resources []int `json:"resources"`
	Gems      []int `json:"gems"`
}

func (conn *Conn) Send(s *GameState) error {
	state := *s
	stateJSON := GameStateJSON{
		Type:       "state",
		TeamID:     conn.TeamID,
		UserID:     int(conn.UserID),
		LockHolder: make([]Queue, len(state.Teams)),
		Score:      make([]Score, len(state.Teams)),
	}
	for i, team := range state.Teams {
		stateJSON.LockHolder[i] = team.LockHolder
		stateJSON.Score[i] = Score{
			Resources: team.ResourceBalance,
			Gems:      team.GemBalance,
		}
	}
	b, err := json.Marshal(stateJSON)
	if err != nil {
		return fmt.Errorf("error sending message: cannot marshal message: %v", err)
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()
	return conn.ws.WriteMessage(websocket.TextMessage, b)
}
