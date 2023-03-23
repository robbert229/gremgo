package gremgo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Client is a container for the gremgo client.
type Client struct {
	conn             dialer
	requests         chan []byte
	responses        chan []byte
	results          *sync.Map
	resultsErr       *sync.Map
	responseNotifier *sync.Map // responseNotifier notifies the requester that a response has arrived for the request
	respMutex        *sync.Mutex
	Errored          bool
	logger           Logger
}

// NewDialer returns a WebSocket dialer to use when connecting to Gremlin Server
func NewDialer(host string, configs ...DialerConfig) (dialer *Ws) {
	dialer = &Ws{
		timeout:      5 * time.Second,
		pingInterval: 60 * time.Second,
		writingWait:  15 * time.Second,
		readingWait:  15 * time.Second,
		connected:    false,
		quit:         make(chan struct{}),
		logger:       NopLogger{},
	}

	for _, conf := range configs {
		conf(dialer)
	}

	dialer.host = host
	return dialer
}

func newClient() (c Client) {
	c.requests = make(chan []byte, 3)  // c.requests takes any request and delivers it to the WriteWorker for dispatch to Gremlin Server
	c.responses = make(chan []byte, 3) // c.responses takes raw responses from ReadWorker and delivers it for sorting to handelResponse
	c.results = &sync.Map{}
	c.resultsErr = &sync.Map{}
	c.responseNotifier = &sync.Map{}
	c.respMutex = &sync.Mutex{} // c.mutex ensures that sorting is thread safe
	c.logger = NopLogger{}
	return
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
//
// Deprecated: use DialContext instead.
func Dial(conn dialer, errs chan error) (c Client, err error) {
	return DialContext(context.Background(), conn, errs)
}

// DialContext returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func DialContext(ctx context.Context, conn dialer, errs chan error) (c Client, err error) {
	c = newClient()
	c.conn = conn

	// Connects to Gremlin Server
	err = conn.connect(ctx)
	if err != nil {
		return
	}

	quit := conn.(*Ws).quit

	go c.writeWorker(errs, quit)
	go c.readWorker(errs, quit)
	go conn.ping(errs)

	return
}

func (c *Client) executeRequest(ctx context.Context, query string, bindings, rebindings map[string]string) (interface{}, error) {
	req, err := prepareRequest(query, bindings, rebindings)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	msg, err := packageRequest(req)
	if err != nil {
		// TODO(robbert229) this Log call is here because there was originally
		// one logging directly to stdout/stderr here. We could remove it, but
		// not doing so currently.
		c.logger.Log("err", err, "msg", "failed to package request")

		return nil, fmt.Errorf("failed to package request: %w", err)
	}

	c.responseNotifier.Store(req.RequestId, make(chan int, 1))
	c.dispatchRequest(msg)

	resp, err := c.retrieveResponse(req.RequestId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve response: %w", err)
	}

	return resp, nil
}

func (c *Client) authenticate(requestId string) (error) {
	auth := c.conn.getAuth()
	req, err := prepareAuthRequest(requestId, auth.username, auth.password)
	if err != nil {
		return fmt.Errorf("failed to prepare auth request: %w", err)
	}

	msg, err := packageRequest(req)
	if err != nil {
		// TODO(robbert229) this Log call is here because there was originally
		// one logging directly to stdout/stderr here. We could remove it, but
		// not doing so currently.
		c.logger.Log("err", err, "msg", "failed to package auth request")
		
		return fmt.Errorf("failed to package auth request: %w", err)
	}

	c.dispatchRequest(msg)

	return nil
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
//
// Deprecated: use ExecuteContext instead.
func (c *Client) Execute(query string, bindings, rebindings map[string]string) (interface{}, error) {
	return c.ExecuteContext(context.Background(), query, bindings, rebindings)
}

// ExecuteContext formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteContext(ctx context.Context, query string, bindings, rebindings map[string]string) (resp interface{}, err error) {
	if c.conn.isDisposed() {
		return nil, errors.New("you cannot write on disposed connection")
	}
	resp, err = c.executeRequest(ctx, query, bindings, rebindings)
	return
}

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
//
// Deprecated: use ExecuteFileContext instead.
func (c *Client) ExecuteFile(path string, bindings, rebindings map[string]string) (resp interface{}, err error) {
	return c.ExecuteFileContext(context.Background(), path, bindings, rebindings)
}

// ExecuteFileContext takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteFileContext(ctx context.Context, path string, bindings, rebindings map[string]string) (interface{}, error) {
	if c.conn.isDisposed() {
		return nil, errors.New("you cannot write on disposed connection")
	}

	d, err := os.ReadFile(path) // Read script from file
	if err != nil {
		// TODO(robbert229) this Log call is here because there was originally
		// one logging directly to stdout/stderr here. We could remove it, but
		// not doing so currently.
		c.logger.Log("err", err, "msg", "failed to read file")

		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	query := string(d)

	resp, err := c.executeRequest(ctx, query, bindings, rebindings)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Close closes the underlying connection and marks the client as closed.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.close()
	}
}
