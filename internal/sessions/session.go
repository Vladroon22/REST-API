package sessions

import (
	"sync"
	"time"
)

type MapSessions struct {
	session map[int]*Session
	mu      sync.Mutex
}

type Session struct {
	ID         int
	regTime    time.Time
	expireTime time.Time
	token      string
}

func NewSession() *MapSessions {
	return &MapSessions{
		session: make(map[int]*Session),
	}
}

func (s *MapSessions) AddNewSeesion(token string, userID int, dur time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session[userID] = &Session{
		ID:      userID,
		token:   token,
		regTime: time.Now(),
	}

	go s.CheckSession(token, userID, dur)
}

func (s *MapSessions) CheckSession(token string, id int, dur time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	go expiresTime(dur)

	if session, exist := s.session[id]; exist {
		if session.expireTime.Before(session.regTime) {
			s.DeleteSession(id)
		}
	}
}

func expiresTime(d time.Duration) {
	tk := time.NewTimer(d)
	<-tk.C
}

func (s *MapSessions) DeleteSession(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.session, id)
}
