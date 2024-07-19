package main

import (
	"fmt"
	"log"
	"time"
)

type Command interface {
	cmd()
}

type Auth struct {
	Nickname string `json:"nickname"`
	Secret   string `json:"secret"`
}

type Lock struct {
	nickname string
	teamID   int
}

type Notice struct {
	Timestamp int    `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Dismissed bool   `json:"dismissed,omitempty"`
}

type NoticeStatusUpdate struct {
	ID        int  `json:"id"`
	Dismissed bool `json:"dismissed,omitempty"`
}

type Transfer struct {
	From   int            `json:"from"`
	To     int            `json:"to"`
	Amount map[string]int `json:"amount"`
}

// timeTick is sent every second to update game state.
type timeTick struct{}

func (Auth) cmd()               {}
func (Lock) cmd()               {}
func (Notice) cmd()             {}
func (NoticeStatusUpdate) cmd() {}
func (Transfer) cmd()           {}
func (timeTick) cmd()           {}

type GameState struct {
	Players []*Conn
	Teams   []Team
	Notices []Notice

	msg chan<- Command
	cfg Config
}

type Team struct {
	TeamConfig
	Balance              []int  `json:"balance"`
	LockHolder           string `json:"lock_holder,omitempty"`
	LockSecondsRemaining int    `json:"lock_seconds_remaining,omitempty"`
}

func NewGameState(cfg Config) (*GameState, error) {
	teams := make([]Team, 0, len(cfg.Teams))
	for i, team := range cfg.Teams {
		balance := make([]int, len(cfg.Currencies))
		switch {
		case team.InitialBalance == nil:
			// Default everything to zero if not specified.
		case len(*team.InitialBalance) != len(cfg.Currencies):
			return nil, fmt.Errorf(
				"error validating config: expect initial balance to be of length %d for team %d (%s)",
				len(cfg.Currencies), i, team.Name,
			)
		}
		teams = append(teams, Team{
			TeamConfig: team,
			Balance:    balance,
		})
	}
	return &GameState{
		Players: make([]*Conn, 0, 10),
		Notices: make([]Notice, 0, 10),
		Teams:   teams,
		cfg:     cfg,
	}, nil
}

func (state *GameState) Start() {
	hasUpdate := false
	ch := make(chan Command)
	state.msg = ch
	go func() {
		for range time.NewTicker(time.Second).C {
			ch <- timeTick{}
		}
	}()
	msg := <-ch
	for {
		state.handleMessage(msg)
		var ok bool
		msg, ok = <-ch
		if ok {
			// Process the message.
			continue
		}
		// Probably processed all the message. Time to send updates.
		state.gc()
		if hasUpdate {
			state.send()
		}
		msg = <-ch
	}
}

func (state *GameState) handleMessage(cmd Command) {
	switch cmd := cmd.(type) {
	case Auth:
		panic("unexpected Auth message in GameState")
	case Lock:
		state.Teams[cmd.teamID].LockHolder = cmd.nickname
		state.Teams[cmd.teamID].LockSecondsRemaining = state.cfg.LockLengthSeconds
	case Notice:
		cmd.Timestamp = int(time.Now().Unix())
		state.Notices = append(state.Notices, cmd)
	case NoticeStatusUpdate:
		if cmd.ID >= len(state.Notices) {
			log.Printf(
				"unexpected malformed NoticeStatusUpdate: ID %d for length %d\n",
				cmd.ID, len(state.Notices),
			)
			return
		}
		state.Notices[cmd.ID].Dismissed = cmd.Dismissed
	default:
		panic(fmt.Sprintf("unknown command type: %T", cmd))
	}
}

func (state *GameState) send() {
	for _, v := range state.Players {
		v := v
		// FIXME: This is racey.
		go func() {
			if err := v.Send(state); err != nil {
				log.Printf("failed to send game state to connection: %v\n", err)
				v.Close()
			}
		}()
	}
}

func (state *GameState) gc() {
	// Compress Players to only keep those that are still connected (not closed).
	write := 0
	for i := 0; i < len(state.Players); i++ {
		if state.Players[i].closed {
			continue
		}
		state.Players[write] = state.Players[i]
		write++
	}
	state.Players = state.Players[:write]
}
