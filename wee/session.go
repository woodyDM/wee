package wee

import (
	log "log"
	"net/http"
	"strconv"
	"time"
)

const (
	GSession = "gsession"
)

type Session interface {
	SessionId() string
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}
type SessionIdGenerator func() string

type SessionHolder interface {
	GetSession(sessionId string) (Session, bool)
	SetSession(sessionId string, s Session)
}

type SessionMidWare struct {
	ExpireSeconds int
	holder        SessionHolder
	IdGenerator   SessionIdGenerator
}

type mapSessionHolder struct {
	content map[string]Session
}

func NewSessionMidware(expireSeconds int, holder SessionHolder, g SessionIdGenerator) *SessionMidWare {
	return &SessionMidWare{
		ExpireSeconds: expireSeconds,
		holder:        holder,
		IdGenerator:   g,
	}
}

func NewSimpleSessionMidware(expireSeconds int, holder SessionHolder) *SessionMidWare {
	return NewSessionMidware(expireSeconds, holder, generateSessionKey)
}

func NewMapSessionMidware() *SessionMidWare {
	return NewSessionMidware(3600, &mapSessionHolder{
		content: map[string]Session{},
	}, generateSessionKey)
}

func (m *mapSessionHolder) GetSession(key string) (Session, bool) {
	session := m.content[key]
	return session, session != nil

}

func (m *mapSessionHolder) SetSession(key string, s Session) {
	m.content[key] = s
}

type mapSession struct {
	id      string
	content map[string]string
}

func (m *mapSession) Delete(key string) {
	delete(m.content, key)
}

func (m *mapSession) SessionId() string {
	return m.id
}

func (m *mapSession) Get(key string) (string, bool) {
	v, ok := m.content[key]
	return v, ok
}

func (m *mapSession) Set(key string, value string) {
	m.content[key] = value
}

func (s *SessionMidWare) Action(ctx *Context, chain *MidWareChain) {
	sessionKey, ok := getCurrentSessionKey(ctx)
	if !ok {
		sessionKey = writeSessionKeyCookie(ctx, s.ExpireSeconds)
	}
	s.bindSession(sessionKey, ctx)
	chain.Next(ctx)
}

func (s *SessionMidWare) bindSession(sessionId string, ctx *Context) {
	session, ok := s.holder.GetSession(sessionId)
	if !ok {
		session = &mapSession{
			content: make(map[string]string),
			id:      sessionId,
		}
		s.holder.SetSession(sessionId, session)
	}
	ctx.Session = session
}

func writeSessionKeyCookie(ctx *Context, expire int) string {
	key := generateSessionKey()
	cookie := http.Cookie{
		Value:    key,
		Name:     GSession,
		HttpOnly: true,
		MaxAge:   expire,
		Path:     "/",
	}
	ctx.Response.Header().Set("Set-Cookie", cookie.String())
	return key
}

func generateSessionKey() string {
	key := time.Now().Nanosecond()
	return strconv.Itoa(key)
}

func getCurrentSessionKey(ctx *Context) (string, bool) {

	if cookie, e := ctx.Request.Cookie(GSession); e == nil {
		sessionKey := cookie.Value
		if sessionKey != "" {
			log.Printf("当前 session %s\n", sessionKey)
			return sessionKey, true
		}
	}
	log.Println("no session key found")
	return "", false
}
