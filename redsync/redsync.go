package redsync

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

var errLockFailed = errors.New("redsync: failed to lock")
var errUnlockFailed = errors.New("redsync: failed to unlock")

// Mutex Mutex
type Mutex struct {
	name   string
	pool   *redis.Pool
	delay  time.Duration
	value  string
	expiry time.Duration
}

// NewMutex NewMutex
func NewMutex(name string, pool *redis.Pool) *Mutex {
	return &Mutex{
		name:   name,
		pool:   pool,
		delay:  200 * time.Millisecond,
		expiry: 6 * time.Second,
	}
}

// Lock Lock
func (m *Mutex) Lock() error {
	value, err := m.genValue()
	if err != nil {
		return err
	}
	for i := 0; i < 60; i++ {
		if i != 0 {
			time.Sleep(m.delay)
		}
		conn := m.pool.Get()
		reply, err := redis.String(conn.Do("SET", m.name, value, "NX", "PX", int(m.expiry/time.Millisecond)))
		conn.Close()
		if err == nil && reply == "OK" {
			m.value = value
			return nil
		}
	}
	return errLockFailed
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

// Unlock Unlock
func (m *Mutex) Unlock() error {
	conn := m.pool.Get()
	status, err := deleteScript.Do(conn, m.name, m.value)
	conn.Close()
	if err != nil {
		return err
	}
	if status == 0 {
		return errUnlockFailed
	}
	return nil
}

func (m *Mutex) genValue() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
