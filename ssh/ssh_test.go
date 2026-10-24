package ssh

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func testSigner(t *testing.T) (gossh.Signer, string) {
	t.Helper()
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := gossh.NewSignerFromKey(key)
	if err != nil {
		t.Fatal(err)
	}
	data, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	filename := filepath.Join(t.TempDir(), "id_ed25519")
	if err := os.WriteFile(filename, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: data}), 0600); err != nil {
		t.Fatal(err)
	}
	return signer, filename
}

type localSSHServer struct {
	address          string
	knownHosts       string
	keyfile          string
	done             chan struct{}
	sessionRequested chan struct{}
}

func startSSHServer(t *testing.T, network, address string, rejectSession, stallSession bool) *localSSHServer {
	t.Helper()
	listener, err := net.Listen(network, address)
	if err != nil {
		if network == "tcp6" {
			t.Skipf("IPv6 loopback unavailable: %v", err)
		}
		t.Fatal(err)
	}
	hostSigner, _ := testSigner(t)
	clientSigner, clientKeyfile := testSigner(t)
	knownHosts := filepath.Join(t.TempDir(), "known_hosts")
	if err := os.WriteFile(knownHosts, []byte(knownhosts.Line([]string{listener.Addr().String()}, hostSigner.PublicKey())+"\n"), 0600); err != nil {
		listener.Close()
		t.Fatal(err)
	}
	server := &localSSHServer{
		address: listener.Addr().String(), knownHosts: knownHosts, keyfile: clientKeyfile,
		done: make(chan struct{}), sessionRequested: make(chan struct{}),
	}
	config := &gossh.ServerConfig{PublicKeyCallback: func(_ gossh.ConnMetadata, offered gossh.PublicKey) (*gossh.Permissions, error) {
		if !bytes.Equal(offered.Marshal(), clientSigner.PublicKey().Marshal()) {
			return nil, errors.New("unexpected client key")
		}
		return nil, nil
	}}
	config.AddHostKey(hostSigner)
	var accepted net.Conn
	var lock sync.Mutex
	var requested sync.Once
	go func() {
		defer close(server.done)
		transport, err := listener.Accept()
		if err != nil {
			return
		}
		lock.Lock()
		accepted = transport
		lock.Unlock()
		defer transport.Close()
		connection, channels, requests, err := gossh.NewServerConn(transport, config)
		if err != nil {
			return
		}
		defer connection.Close()
		go gossh.DiscardRequests(requests)
		for incoming := range channels {
			if incoming.ChannelType() != "session" {
				_ = incoming.Reject(gossh.UnknownChannelType, "sessions only")
				continue
			}
			requested.Do(func() { close(server.sessionRequested) })
			if stallSession {
				continue // Neither accept nor reject: exercise the session-creation deadline.
			}
			if rejectSession {
				_ = incoming.Reject(gossh.Prohibited, "sessions disabled")
				continue
			}
			channel, requests, err := incoming.Accept()
			if err != nil {
				return
			}
			go func() {
				defer channel.Close()
				for request := range requests {
					if request.Type != "exec" {
						_ = request.Reply(false, nil)
						continue
					}
					_ = request.Reply(true, nil)
					_, _ = io.WriteString(channel, "hello from SSH\n")
					_, _ = channel.SendRequest("exit-status", false, gossh.Marshal(struct{ Status uint32 }{0}))
					return
				}
			}()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		lock.Lock()
		if accepted != nil {
			_ = accepted.Close()
		}
		lock.Unlock()
		select {
		case <-server.done:
		case <-time.After(3 * time.Second):
			t.Error("local SSH server did not stop")
		}
	})
	return server
}

func waitForSSHClose(t *testing.T, server *localSSHServer) {
	t.Helper()
	select {
	case <-server.done:
	case <-time.After(3 * time.Second):
		t.Fatal("SSH connection was not released")
	}
}

func TestSecureShellWithHostKey(t *testing.T) {
	for _, address := range []string{"127.0.0.1:0", "[::1]:0"} {
		t.Run(address, func(t *testing.T) {
			network := "tcp"
			if address == "[::1]:0" {
				network = "tcp6"
			}
			server := startSSHServer(t, network, address, false, false)
			host, port, err := net.SplitHostPort(server.address)
			if err != nil {
				t.Fatal(err)
			}
			client, session, err := SecureShellWithHostKey("test", host, port, server.keyfile, server.knownHosts)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			defer session.Close()
			output, err := session.Output("hello")
			if err != nil || string(output) != "hello from SSH\n" {
				t.Fatalf("Output = %q, %v", output, err)
			}
		})
	}
}

func TestLegacySessionRemainsUsableAndOwnsConnection(t *testing.T) {
	server := startSSHServer(t, "tcp", "127.0.0.1:0", false, false)
	host, port, _ := net.SplitHostPort(server.address)
	session := SecureShell("test", host, port, server.keyfile)
	defer session.Close()
	output, err := session.Output("hello")
	if err != nil || string(output) != "hello from SSH\n" {
		t.Fatalf("returned session was already closed: %q, %v", output, err)
	}
	_ = session.Close()
	waitForSSHClose(t, server)
}

func TestVerifiedHostKeyRejectsOtherServer(t *testing.T) {
	server := startSSHServer(t, "tcp", "127.0.0.1:0", false, false)
	other, _ := testSigner(t)
	if err := os.WriteFile(server.knownHosts, []byte(knownhosts.Line([]string{server.address}, other.PublicKey())+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	host, port, _ := net.SplitHostPort(server.address)
	client, session, err := SecureShellWithHostKey("test", host, port, server.keyfile, server.knownHosts)
	if err == nil || client != nil || session != nil {
		t.Fatalf("accepted mismatched host key: %v, %v, %v", client, session, err)
	}
	waitForSSHClose(t, server)
}

func TestVerifiedKeyAndKnownHostsErrors(t *testing.T) {
	server := startSSHServer(t, "tcp", "127.0.0.1:0", false, false)
	missing := filepath.Join(t.TempDir(), "missing")
	for _, paths := range [][2]string{{server.keyfile, missing}, {missing, server.knownHosts}} {
		client, session, err := SecureShellWithHostKey("test", "127.0.0.1", "22", paths[0], paths[1])
		if err == nil || client != nil || session != nil {
			t.Fatal("invalid credential paths did not return an error")
		}
	}
	if err := os.WriteFile(server.keyfile, []byte("not a private key"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := SecureShellWithHostKey("test", "127.0.0.1", "22", server.keyfile, server.knownHosts); err == nil {
		t.Fatal("invalid private key was accepted")
	}
	if _, _, err := connectToHostWithPublickey("test", "127.0.0.1:22", missing); err == nil {
		t.Fatal("legacy path accepted a missing private key")
	}
}

func TestSessionCreationFailureClosesConnection(t *testing.T) {
	server := startSSHServer(t, "tcp", "127.0.0.1:0", true, false)
	host, port, _ := net.SplitHostPort(server.address)
	client, session, err := SecureShellWithHostKey("test", host, port, server.keyfile, server.knownHosts)
	if err == nil || client != nil || session != nil {
		t.Fatal("rejected session did not return an error")
	}
	waitForSSHClose(t, server)
}

func TestConnectionSetupDeadlineIncludesSessionCreation(t *testing.T) {
	server := startSSHServer(t, "tcp", "127.0.0.1:0", false, true)
	signer, err := readSigner(server.keyfile)
	if err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() {
		_, _, err := dialSession(server.address, &gossh.ClientConfig{
			User: "test", Auth: []gossh.AuthMethod{gossh.PublicKeys(signer)}, HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		}, 200*time.Millisecond, false)
		result <- err
	}()
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("stalled session creation succeeded")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("session creation did not respect the setup deadline")
	}
	select {
	case <-server.sessionRequested:
	default:
		t.Fatal("test never reached session creation")
	}
	waitForSSHClose(t, server)
}

func TestHandshakeDeadlineAndDialFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		connection, err := listener.Accept()
		if err == nil {
			defer connection.Close()
			_, _ = io.Copy(io.Discard, connection) // Read the client banner, but send none.
		}
	}()
	result := make(chan error, 1)
	go func() {
		_, _, err := dialSession(listener.Addr().String(), &gossh.ClientConfig{HostKeyCallback: gossh.InsecureIgnoreHostKey()}, 50*time.Millisecond, false)
		result <- err
	}()
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("unresponsive peer did not time out")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("SSH handshake did not time out")
	}
	select {
	case <-closed:
	case <-time.After(3 * time.Second):
		t.Fatal("timed-out handshake leaked a connection")
	}
	if _, _, err := dialSession("invalid address", &gossh.ClientConfig{}, time.Second, false); err == nil {
		t.Fatal("invalid dial address was accepted")
	}
}

type countingConn struct {
	gossh.Conn
	closes int32
}

func (connection *countingConn) Close() error {
	atomic.AddInt32(&connection.closes, 1)
	return nil
}

func TestOwnedSessionCloseIsConcurrentAndIdempotent(t *testing.T) {
	owner := &countingConn{}
	channel := &owningSessionChannel{owner: owner}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := channel.Close(); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if count := atomic.LoadInt32(&owner.closes); count != 1 {
		t.Fatalf("connection was closed %d times", count)
	}
}
