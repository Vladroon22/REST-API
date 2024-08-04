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
	ID    int
	d     time.Time
	token string
}

func NewSession() *MapSessions {
	return &MapSessions{
		session: make(map[int]*Session),
	}
}

func (s *MapSessions) AddNewSeesion(token string, userID int, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session[userID] = &Session{
		ID:    userID,
		token: token,
	}

	go s.CheckSession(token, userID, d)
}

func (s *MapSessions) CheckSession(token string, id int, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	go expiresTime(d)

	if session, exist := s.session[id]; !exist {
		if session.d.Before(time.Now()) {
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
