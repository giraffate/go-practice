package mem

import (
	"errors"
	"time"

	"../../session"
)

//SessionStore stores session.
type SessionStore struct {
	sid       string
	m         map[string]interface{}
	createdAt time.Time
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
}

// SessionCreate initiates session.
func (pder *Provider) SessionCreate(sid string) (session.Session, error) {
	ss := SessionStore{sid: sid, m: make(map[string]interface{}), createdAt: time.Now()}
	pder.sessions[sid] = ss
	return &ss, nil
}

// SessionRead return session.
func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	ss, ok := pder.sessions[sid]
	if ok {
		return &ss, nil
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

func init() {
	sessions := make(map[string]SessionStore)
	session.Register("mem", &Provider{sessions: sessions})
}
