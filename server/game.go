package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Command interface {
	cmd()
}

type ConnCommand struct {
	Conn    *Conn
	Command Command
}

type Auth struct {
	Nickname string `json:"nickname"`
	Secret   string `json:"secret"`
}

type Lock struct{}

type Unlock struct{}

type NoticeLevel string

const (
	NoticePause     NoticeLevel = "pause"
	NoticeHighlight NoticeLevel = "highlight"
	NoticeWarning   NoticeLevel = "warning"
	NoticeMessage   NoticeLevel = "message"
)

type Notice struct {
	ID        int         `json:"id"`
	Timestamp int         `json:"timestamp"`
	Level     NoticeLevel `json:"level"`
	Dismissed bool        `json:"dismissed,omitempty"`
	// TeamID is the intended recipient of the notice. If not specified, it
	// will be sent to everyone.
	TeamID *int `json:"team_id,omitempty"`

	// Message or (TranslationKey and TranslationValue) will be present.
	Message          string            `json:"message,omitempty"`
	TranslationKey   string            `json:"translation_key,omitempty"`
	TranslationValue map[string]string `json:"translation_value,omitempty"`

	// skipAdmin is set to true if the notice does not need to be sent to the
	// admins.
	skipAdmin bool
}

// Sync with frontend: src/lib/{en,zh}.json.
const (
	TranslationTransferSent     = "dashboard.label.transferSent"
	TranslationTransferReceived = "dashboard.label.transferReceived"
)

type NoticeStatusUpdate struct {
	ID        int  `json:"id"`
	Dismissed bool `json:"dismissed,omitempty"`
}

type Transfer struct {
	From           int   `json:"from"`
	To             int   `json:"to"`
	GemAmount      []int `json:"gem_amount"`
	ResourceAmount []int `json:"resource_amount"`
}

// timeTick is sent every second to update game state.
type timeTick struct{}

func (Auth) cmd()               {}
func (Lock) cmd()               {}
func (Unlock) cmd()             {}
func (Notice) cmd()             {}
func (NoticeStatusUpdate) cmd() {}
func (Transfer) cmd()           {}
func (timeTick) cmd()           {}

type GameState struct {
	Players []*Conn
	Teams   []Team
	Notices []Notice

	msg          chan<- ConnCommand
	cfg          Config
	nextNoticeID int
}

type Team struct {
	TeamConfig
	ResourceBalance      []int `json:"resource_balance"`
	GemBalance           []int `json:"gem_balance"`
	LockHolder           Queue `json:"lock_holder,omitempty"`
	LockSecondsRemaining int   `json:"lock_seconds_remaining,omitempty"`
}

func NewGameState(cfg Config) (*GameState, error) {
	teams := make([]Team, 0, len(cfg.Teams))
	for i, team := range cfg.Teams {
		resourceBal := make([]int, len(cfg.ResourceNames))
		gemBal := make([]int, len(cfg.GemNames))
		switch {
		case team.InitialBalance == nil:
			// Default everything to zero if not specified.
		case len(team.InitialBalance.Resources) != len(cfg.ResourceNames):
			return nil, fmt.Errorf(
				"error validating config: expect initial resource balance to be of length %d for team %d (%s)",
				len(cfg.ResourceNames), i, team.Name,
			)
		case len(team.InitialBalance.Gems) != len(cfg.GemNames):
			return nil, fmt.Errorf(
				"error validating config: expect initial gem balance to be of length %d for team %d (%s)",
				len(cfg.GemNames), i, team.Name,
			)
		default:
			copy(resourceBal, team.InitialBalance.Resources)
			copy(gemBal, team.InitialBalance.Gems)
		}
		if team.Admin {
			// Admins do not have a balance.
			resourceBal = nil
			gemBal = nil
		}
		teams = append(teams, Team{
			TeamConfig:      team,
			ResourceBalance: resourceBal,
			GemBalance:      gemBal,
			LockHolder:      NewQueue(cfg.LockLengthSeconds),
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
	ch := make(chan ConnCommand)
	state.msg = ch
	go func() {
		for range time.NewTicker(time.Second).C {
			ch <- ConnCommand{
				Conn:    nil,
				Command: timeTick{},
			}
		}
	}()
	msg := <-ch
	for {
		state.handleMessage(msg)
		select {
		case msg = <-ch:
			// Process the message.
			continue
		default:
			// Probably processed all the message. Time to send updates.
			state.gc()
			state.sendState()
			// Block until next message.
			msg = <-ch
		}
	}
}

func (state *GameState) handleMessage(connCmd ConnCommand) {
	conn := connCmd.Conn
	if _, isTimeTick := connCmd.Command.(timeTick); conn == nil && !isTimeTick {
		panic("expected conn to be filled in when processed in GameState")
	}
	switch cmd := connCmd.Command.(type) {
	case Auth:
		state.Players = append(state.Players, conn)
	case Lock:
		state.Teams[conn.TeamID].LockHolder.Join(NewLockHolder(connCmd.Conn))
		/*
			state.Teams[conn.UserID].LockSecondsRemaining = state.cfg.LockLengthSeconds
		*/
	case Unlock:
		state.Teams[conn.TeamID].LockHolder.Leave(NewLockHolder(connCmd.Conn))
	case Notice:
		cmd.Timestamp = int(time.Now().Unix())
		cmd.ID = state.nextNoticeID
		state.nextNoticeID++
		switch cmd.Level {
		case NoticePause, NoticeHighlight:
			// Pause and highlights are important notices that we want to
			// re-send to players whenever they reconnect.
			state.Notices = append(state.Notices, cmd)
		case NoticeWarning, NoticeMessage:
			// These are OK to drop as they are typically one-off.
		default:
			panic(fmt.Sprintf("unknown notice level: %s", cmd.Level))
		}
		state.sendNotice(cmd)
	case NoticeStatusUpdate:
		for i, notice := range state.Notices {
			if notice.ID != cmd.ID {
				continue
			}
			state.Notices[i].Dismissed = cmd.Dismissed
			state.sendNotice(state.Notices[i])
			return
		}
		log.Printf(
			"unexpected malformed NoticeStatusUpdate: ID %d for length %d\n",
			cmd.ID, len(state.Notices),
		)
	case Transfer:
		state.handleTransferCommand(connCmd)
	case timeTick:
		for i := range state.Teams {
			state.Teams[i].LockHolder.Advance(1)
		}
	default:
		panic(fmt.Sprintf("unknown command type: %T", cmd))
	}
}

func (state *GameState) sendState() {
	for _, v := range state.Players {
		v := v
		// FIXME: This is racey.
		go func() {
			if err := v.SendState(state, false); err != nil {
				log.Printf("failed to send game state to connection: %v\n", err)
				v.Close()
			}
		}()
	}
}

func (state *GameState) sendNotice(notice Notice) {
	for _, v := range state.Players {
		v := v
		if !notice.AppliesTo(v, state) {
			continue
		}
		// FIXME: This is racey.
		go func() {
			if err := v.SendNotice(notice); err != nil {
				log.Printf("failed to send notice to connection: %v\n", err)
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

func (state *GameState) handleTransferCommand(transfer ConnCommand) {
	cmd := transfer.Command.(Transfer)
	// Check From.
	switch {
	case state.cfg.Teams[cmd.From].Admin:
		// Admins are always allowed.
	default:
		fromTeam := state.Teams[cmd.From]
		for i := range cmd.GemAmount {
			if cmd.GemAmount[i] > 0 && fromTeam.GemBalance[i] < cmd.GemAmount[i] {
				log.Printf("Insufficient gems to transfer.")
				return
			}
		}
		for i := range cmd.ResourceAmount {
			if cmd.ResourceAmount[i] > 0 && fromTeam.ResourceBalance[i] < cmd.ResourceAmount[i] {
				log.Printf("Insufficient resources to transfer.")
				return
			}
		}
		// Deduct from players.
		for i := range cmd.GemAmount {
			fromTeam.GemBalance[i] -= cmd.GemAmount[i]
		}
		for i := range cmd.ResourceAmount {
			fromTeam.ResourceBalance[i] -= cmd.ResourceAmount[i]
		}
	}
	// Update To.
	switch {
	case state.cfg.Teams[cmd.To].Admin:
		// Admins' balance does not need to be updated.
	default:
		toTeam := state.Teams[cmd.To]
		for i := range cmd.GemAmount {
			toTeam.GemBalance[i] += cmd.GemAmount[i]
		}
		for i := range cmd.ResourceAmount {
			toTeam.ResourceBalance[i] += cmd.ResourceAmount[i]
		}
	}
	// Send notices.
	noticeID := state.nextNoticeID
	state.nextNoticeID++
	resources := make([]string, 0, len(state.cfg.ResourceNames)+len(state.cfg.GemNames))
	for i, v := range cmd.GemAmount {
		if v != 0 {
			resources = append(resources, fmt.Sprintf("%dx %s", v, state.cfg.GemNames[i]))
		}
	}
	for i, v := range cmd.ResourceAmount {
		if v != 0 {
			resources = append(resources, fmt.Sprintf("%dx %s", v, state.cfg.ResourceNames[i]))
		}
	}
	resourceCombined := strings.Join(resources, ", ")
	state.sendNotice(Notice{
		ID:             noticeID,
		Timestamp:      int(time.Now().Unix()),
		Level:          NoticeMessage,
		TeamID:         &cmd.From,
		TranslationKey: TranslationTransferSent,
		TranslationValue: map[string]string{
			"playerSender": transfer.Conn.Nickname,
			"recipient":    state.cfg.Teams[cmd.To].Name,
			"resource":     resourceCombined,
		},
		skipAdmin: true,
	})
	state.sendNotice(Notice{
		ID:             noticeID,
		Timestamp:      int(time.Now().Unix()),
		Level:          NoticeMessage,
		TeamID:         &cmd.To,
		TranslationKey: TranslationTransferReceived,
		TranslationValue: map[string]string{
			"playerSender": transfer.Conn.Nickname,
			"sender":       state.cfg.Teams[cmd.From].Name,
			"resource":     resourceCombined,
		},
		skipAdmin: true,
	})
}

func (notice Notice) AppliesTo(conn *Conn, s *GameState) bool {
	// The player is an admin and this notice is not skipped for admin.
	if s.cfg.Teams[conn.TeamID].Admin && !notice.skipAdmin {
		return true
	}
	// The notice belongs to the team.
	return notice.TeamID == nil || *notice.TeamID == conn.TeamID
}
