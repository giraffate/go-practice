package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	defaultCookieName = "SESSIONID"
)

// Manager is the session manager.
//
// TODO
// Add test
type Manager struct {
	m          sync.Mutex
	provider   Provider
	cookieName string
	maxAge     int // sec
}

// Provider is the interface of session provider.
type Provider interface {
	SessionCreate(sid string, maxAge int) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC() error
}

// Session is the interface of session.
type Session interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Destroy(key string) error
	SessionID() string
}

var providerList = make(map[string]Provider)

// NewManager is a constructor of manager.
func NewManager(providerName string, maxAge int) (*Manager, error) {
	pder, ok := providerList[providerName]
	if !ok {
		return nil, errors.New("No provider is found")
	}
	return &Manager{provider: pder, cookieName: defaultCookieName, maxAge: maxAge}, nil
}

// Register set provider.
func Register(name string, pder Provider) error {
	if pder == nil {
		return errors.New("No provider is set")
	}
	if _, ok := providerList[name]; ok {
		return errors.New(name + " is already registered")
	}
	providerList[name] = pder
	return nil
}

// SessionCreate creates session.
func (m *Manager) SessionCreate(w http.ResponseWriter) (Session, error) {
	m.m.Lock()
	defer m.m.Unlock()
	sid := sessionID()
	session, err := m.provider.SessionCreate(sid, m.maxAge)
	if err != nil {
		return nil, err
	}
	http.SetCookie(w, &http.Cookie{Name: m.cookieName, Value: sid, MaxAge: m.maxAge})
	return session, nil
}

// SessionRead get session.
func (m *Manager) SessionRead(r *http.Request) (Session, error) {
	var session Session
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return nil, err
	}
	session, err = m.provider.SessionRead(cookie.Value)
	if err != nil {
		return nil, err
	}
	log.Println("Session is founded")
	return session, nil
}

// SessionDestroy delete session by sid.
func (m *Manager) SessionDestroy(r *http.Request) error {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		// Request header has no cookie.
		return nil
	}

	// Request header has cookie.
	m.m.Lock()
	defer m.m.Unlock()
	return m.provider.SessionDestroy(cookie.Value)
}

// SessionGC delete expired session
func (m *Manager) SessionGC() error {
	return m.provider.SessionGC()
}

func sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
