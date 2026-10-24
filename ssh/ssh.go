// Package ssh provides SSH session helpers with optional known-host verification.
package ssh

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const connectionTimeout = 15 * time.Second

// SecureShell connects using a private key, or prompts for a password when
// keyfile is empty. It panics on failure. The caller must close the returned
// session; doing so also closes its owning connection. Connection setup,
// including the SSH handshake, is limited to 15 seconds.
//
// Deprecated: this legacy helper skips host-key verification. Use
// SecureShellWithHostKey to verify the server and receive explicit client ownership.
func SecureShell(user string, host string, port string, keyfile string) *ssh.Session {
	address := net.JoinHostPort(host, port)
	var session *ssh.Session
	var err error
	if keyfile != "" {
		_, session, err = connectToHostWithPublickey(user, address, keyfile)
	} else {
		_, session, err = connectToHost(user, address)
	}
	if err != nil {
		panic(err)
	}
	return session
}

// SecureShellWithHostKey connects with public-key authentication, verifying the
// server against knownHostsFile. Connection setup, authentication and session
// creation have a shared 15-second deadline; active sessions have no deadline.
// host may be an unbracketed IPv6 address. The caller owns both returned objects
// and must close the client as well as the session.
func SecureShellWithHostKey(user, host, port, keyfile, knownHostsFile string) (*ssh.Client, *ssh.Session, error) {
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not load known_hosts file: %w", err)
	}
	signer, err := readSigner(keyfile)
	if err != nil {
		return nil, nil, err
	}
	return dialSession(net.JoinHostPort(host, port), &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
	}, connectionTimeout, false)
}

// connectToHost keeps legacy password authentication without host verification.
func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("SSH-Password: ")
	if _, err := fmt.Scanf("%s\n", &pass); err != nil {
		return nil, nil, err
	}
	return dialSession(host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Legacy API; prefer SecureShellWithHostKey.
	}, connectionTimeout, true)
}

func connectToHostWithPublickey(user, host, publickeyfile string) (*ssh.Client, *ssh.Session, error) {
	signer, err := readSigner(publickeyfile)
	if err != nil {
		return nil, nil, err
	}
	return dialSession(host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Legacy API; prefer SecureShellWithHostKey.
	}, connectionTimeout, true)
}

func readSigner(filename string) (ssh.Signer, error) {
	key, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
}

func dialSession(address string, config *ssh.ClientConfig, timeout time.Duration, closeWithSession bool) (*ssh.Client, *ssh.Session, error) {
	deadline := time.Now().Add(timeout)
	transport, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, nil, err
	}
	if err := transport.SetDeadline(deadline); err != nil {
		_ = transport.Close()
		return nil, nil, err
	}
	connection, channels, requests, err := ssh.NewClientConn(transport, address, config)
	if err != nil {
		_ = transport.Close()
		return nil, nil, err
	}
	if closeWithSession {
		connection = &sessionConnection{Conn: connection}
	}
	client := ssh.NewClient(connection, channels, requests)
	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
		return nil, nil, err
	}
	if err := transport.SetDeadline(time.Time{}); err != nil {
		_ = client.Close()
		return nil, nil, err
	}
	return client, session, nil
}

// The legacy API returns only a Session. Give its session channel ownership of
// the connection so Close releases the client without racing the caller's Wait.
type sessionConnection struct{ ssh.Conn }

func (connection *sessionConnection) OpenChannel(kind string, data []byte) (ssh.Channel, <-chan *ssh.Request, error) {
	channel, requests, err := connection.Conn.OpenChannel(kind, data)
	if err != nil {
		return nil, nil, err
	}
	if kind == "session" {
		channel = &owningSessionChannel{Channel: channel, owner: connection.Conn}
	}
	return channel, requests, nil
}

type owningSessionChannel struct {
	ssh.Channel
	owner    ssh.Conn
	once     sync.Once
	closeErr error
}

func (channel *owningSessionChannel) Close() error {
	// Closing the transport also closes its channel and unblocks pending I/O.
	// Do not wait to send a channel-close packet to an unresponsive peer first.
	channel.once.Do(func() { channel.closeErr = channel.owner.Close() })
	return channel.closeErr
}
