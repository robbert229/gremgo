package gremgo

import "time"

// DialerConfig is the struct for defining configuration for WebSocket dialer
type DialerConfig func(*Ws)

// SetAuthentication sets on dialer credentials for authentication
//
// Deprecated: use WithAuthentication instead.
func SetAuthentication(username string, password string) DialerConfig {
	return func(c *Ws) {
		c.auth = &auth{username: username, password: password}
	}
}

// WithAuthentication sets on dialer credentials for authentication
func WithAuthentication(username string, password string) DialerConfig {
	return func(c *Ws) {
		c.auth = &auth{
			username: username,
			password: password,
		}
	}
}

// SetTimeout sets the dial timeout.
//
// Deprecated: use WithTimeout instead.
func SetTimeout(seconds int) DialerConfig {
	return WithTimeout(time.Second * time.Duration(seconds))
}

// WithTimeout sets the dial timeout.
func WithTimeout(duration time.Duration) DialerConfig {
	return func(c *Ws) {
		c.timeout = duration
	}
}

// SetPingInterval sets the interval of ping sending for know is
// connection is alive and in consequence the client is connected
//
// Deprecated: use WithPingInterval instead.
func SetPingInterval(seconds int) DialerConfig {
	return WithPingInterval(time.Second * time.Duration(seconds))
}

// WithPingInterval sets the interval of ping sending for know is
// connection is alive and in consequence the client is connected
func WithPingInterval(duration time.Duration) DialerConfig {
	return func(c *Ws) {
		c.pingInterval = duration
	}
}

// SetWritingWait sets the time for waiting that writing occur
//
// Deprecated: use WithWritingWait instead.
func SetWritingWait(seconds int) DialerConfig {
	return WithWritingWait(time.Duration(seconds) * time.Second)
}

// WithWritingWait sets the time for waiting that writing occur
func WithWritingWait(duration time.Duration) DialerConfig {
	return func(c *Ws) {
		c.writingWait = duration
	}
}

// SetReadingWait sets the time for waiting that reading occur
//
// Deprecated: use WithReadingWait instead.
func SetReadingWait(seconds int) DialerConfig {
	return WithReadingWait(time.Duration(seconds) * time.Second)
}

// WithReadingWait sets the time for waiting that reading occur
func WithReadingWait(duration time.Duration) DialerConfig {
	return func(c *Ws) {
		c.readingWait = duration
	}
}

// WithLogger sets the logger.
func WithLogger(logger Logger) DialerConfig {
	return func(c *Ws) {
		c.logger = logger
	}
}
