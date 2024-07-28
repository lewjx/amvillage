package main

// Queue is a list of lock holders waiting to acquire the lock.
type Queue struct {
	Players          []LockHolder `json:"players"`
	SecondsRemaining int          `json:"seconds_remaining,omitempty"`

	secondsPerPlayer int
}

type LockHolder struct {
	Nickname string `json:"nickname"`
	UserID   int    `json:"user_id"`
}

func NewQueue(secondsPerPlayer int) Queue {
	return Queue{
		Players:          []LockHolder{},
		secondsPerPlayer: secondsPerPlayer,
	}
}

func NewLockHolder(conn *Conn) LockHolder {
	return LockHolder{
		Nickname: conn.Nickname,
		UserID:   conn.UserID,
	}
}

func (q *Queue) Join(owner LockHolder) {
	for _, v := range q.Players {
		if v == owner {
			return
		}
	}
	q.Players = append(q.Players, owner)
	if len(q.Players) == 1 {
		// They automatically get the lock.
		q.resetCountdown()
	}
}

func (q *Queue) Leave(owner LockHolder) {
	for i, v := range q.Players {
		if v == owner {
			q.Players = append(q.Players[:i], q.Players[i+1:]...)
			if i == 0 {
				q.resetCountdown()
			}
			return
		}
	}
}

func (q *Queue) Advance(seconds int) {
	if len(q.Players) == 0 {
		return
	}
	q.SecondsRemaining -= seconds
	if q.SecondsRemaining <= 0 {
		// Make the first player leave if their lock has expired.
		q.Leave(q.Players[0])
	}
}

func (q Queue) Len() int {
	return len(q.Players)
}

func (q *Queue) resetCountdown() {
	if len(q.Players) > 0 {
		q.SecondsRemaining = q.secondsPerPlayer
	} else {
		// Unset seconds remaining.
		q.SecondsRemaining = 0
	}
}
