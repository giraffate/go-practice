package mem

import (
	"container/list"
	"errors"
	"log"
	"time"

	"../../session"
)

//SessionStore stores session.
type SessionStore struct {
	sid       string
	m         map[string]interface{}
	createdAt time.Time
	maxAge    int
}

// Get returns value by key.
func (ss *SessionStore) Get(key string) (interface{}, error) {
	if v, ok := ss.m[key]; ok {
		return v, nil
	}
	return nil, errors.New("No value is found")
}

// Set set value by key.
func (ss *SessionStore) Set(key string, value interface{}) error {
	ss.m[key] = value
	return nil
}

// Destroy delete by key.
func (ss *SessionStore) Destroy(key string) error {
	delete(ss.m, key)
	return nil
}

// SessionID return sid.
func (ss *SessionStore) SessionID() string {
	return ss.sid
}

// Provider provides the memory session store.
type Provider struct {
	sessions map[string]SessionStore
	gcList   *list.List
}

// SessionCreate initiates session.
func (pder *Provider) SessionCreate(sid string, maxAge int) (session.Session, error) {
	ss := SessionStore{
		sid:       sid,
		m:         make(map[string]interface{}),
		createdAt: time.Now(),
		maxAge:    maxAge,
	}
	pder.sessions[sid] = ss
	pder.gcList.PushBack(&ss)
	return &ss, nil
}

// SessionRead return session.
func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	ss, ok := pder.sessions[sid]
	if ok {
		if ss.createdAt.Add(time.Duration(ss.maxAge) * time.Second).After(time.Now()) {
			return &ss, nil
		}
		pder.SessionGC()
	}
	return nil, errors.New("No session is found")
}

// SessionDestroy delete session.
func (pder *Provider) SessionDestroy(sid string) error {
	if _, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
	}
	return nil
}

// SessionGC delete expired sessions.
func (pder *Provider) SessionGC() error {
	log.Println("Session GC start.")
	for {
		e := pder.gcList.Front()
		if e == nil {
			break
		}
		ss, ok := e.Value.(*SessionStore)
		if !ok {
			pder.gcList.Remove(e)
		}
		if ss.createdAt.Add(time.Duration(ss.maxAge) * time.Second).After(time.Now()) {
			break
		}
		pder.SessionDestroy(ss.sid)
		pder.gcList.Remove(e)
		log.Println("Session GC executed.")
	}
	return nil
}

func init() {
	sessions := make(map[string]SessionStore)
	pder := &Provider{
		sessions: sessions,
		gcList:   list.New(),
	}
	session.Register("mem", pder)
}
