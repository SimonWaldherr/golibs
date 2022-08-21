package ssh

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

func SecureShell(user string, host string, port string, keyfile string) {
	var client *ssh.Client
	var session *ssh.Session
	var err error

	if keyfile != "" {
		client, session, err = connectToHostWithPublickey(user, fmt.Sprintf("%v:%v", host, port), keyfile)
	} else {
		client, session, err = connectToHost(user, fmt.Sprintf("%v:%v", host, port))
	}

	defer client.Close()

	if err != nil {
		panic(err)
	}

	return session
}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("SSH-Password: ")
	fmt.Scanf("%s\n", &pass)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

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

func connectToHostWithPublickey(user, host, publickeyfile string) (*ssh.Client, *ssh.Session, error) {
	key, err := ioutil.ReadFile(publickeyfile)
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
		HostKeyCallback: ssh.HostKeyCallback(func(string, net.Addr, ssh.PublicKey) error { return nil }),
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
