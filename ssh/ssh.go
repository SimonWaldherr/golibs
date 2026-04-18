package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SecureShell connects to a remote host via SSH and returns a session.
// When keyfile is empty, the user is prompted for a password.
// Note: this function uses InsecureIgnoreHostKey; use SecureShellWithHostKey
// for production use where host key verification is required.
func SecureShell(user string, host string, port string, keyfile string) *ssh.Session {
	var client *ssh.Client
	var session *ssh.Session
	var err error

	if keyfile != "" {
		client, session, err = connectToHostWithPublickey(user, fmt.Sprintf("%v:%v", host, port), keyfile)
	} else {
		client, session, err = connectToHost(user, fmt.Sprintf("%v:%v", host, port))
	}

	if err != nil {
		panic(err)
	}

	defer client.Close()

	return session
}

// SecureShellWithHostKey connects to a remote host via SSH using public key
// authentication and verifies the host key against the provided known_hosts file.
func SecureShellWithHostKey(user, host, port, keyfile, knownHostsFile string) (*ssh.Client, *ssh.Session, error) {
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not load known_hosts file: %w", err)
	}

	key, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", host, port), &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
	})
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

// connectToHost connects to a remote host via SSH using password authentication.
// WARNING: uses InsecureIgnoreHostKey which skips host verification.
func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("SSH-Password: ")
	fmt.Scanf("%s\n", &pass)

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

// connectToHostWithPublickey connects to a remote host via SSH using public key authentication.
// WARNING: uses InsecureIgnoreHostKey which skips host verification.
func connectToHostWithPublickey(user, host, publickeyfile string) (*ssh.Client, *ssh.Session, error) {
	key, err := os.ReadFile(publickeyfile)
	if err != nil {
		return nil, nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	client, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
	})
	if client == nil || err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
