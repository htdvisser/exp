package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/credentials"
)

const inProcess = "in-process"

type inProcessAuthInfo struct{}

func (inProcessAuthInfo) AuthType() string { return inProcess }

type inProcessCredentials struct {
	ServerName string
}

func (inProcessCredentials) ClientHandshake(_ context.Context, _ string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, inProcessAuthInfo{}, nil
}

func (inProcessCredentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, inProcessAuthInfo{}, nil
}

func (c inProcessCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: inProcess,
		SecurityVersion:  "master",
		ServerName:       c.ServerName,
	}
}

func (c *inProcessCredentials) Clone() credentials.TransportCredentials { return c }

func (c *inProcessCredentials) OverrideServerName(serverName string) error {
	c.ServerName = serverName
	return nil
}

func newInProcessListener(parent context.Context) *inProcessListener {
	ctx, cancel := context.WithCancel(parent)
	return &inProcessListener{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan net.Conn),
	}
}

type inProcessListener struct {
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan net.Conn
}

func (l inProcessListener) Accept() (net.Conn, error) {
	select {
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	case conn := <-l.ch:
		return conn, nil
	}
}

func (l inProcessListener) Close() error {
	l.cancel()
	return nil
}

type inProcessAddr string

func (inProcessAddr) Network() string  { return inProcess }
func (a inProcessAddr) String() string { return string(a) }

func (l inProcessListener) Addr() net.Addr { return inProcessAddr(inProcess) }

func inProcessDialer(lis *inProcessListener) func(string, time.Duration) (net.Conn, error) {
	return func(addr string, timeout time.Duration) (net.Conn, error) {
		server, client := net.Pipe()
		select {
		case <-time.After(timeout):
			return nil, context.DeadlineExceeded
		case lis.ch <- server:
			return client, nil
		}
	}
}
