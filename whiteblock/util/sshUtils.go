package util

import (
	"errors"
	"fmt"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

type SshClient struct {
	clients []*ssh.Client
}

func NewSshClient(host string) (*SshClient, error) {
	out := new(SshClient)
	for i := 10; i > 0; i -= 5 {
		client, err := sshConnect(host)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		out.clients = append(out.clients, client)
	}
	return out, nil
}

func (this SshClient) GetSession() (*ssh.Session, error) {
	for _, client := range this.clients {
		session, err := client.NewSession()
		if err != nil {
			continue
		}
		return session, nil
	}
	return nil, errors.New("Unable to get a session")
}

/**
 * Easy shorthand for multiple calls to sshExec
 * @param  ...string    commands    The commands to execute
 * @return []string                 The results of the execution of each command
 */
func (this SshClient) MultiRun(commands ...string) ([]string, error) {
	out := []string{}
	for _, command := range commands {
		res, err := this.Run(command)
		if err != nil {
			return nil, err
		}
		out = append(out, string(res))
	}
	return out, nil
}

func (this SshClient) Run(command string) (string, error) {
	session, err := this.GetSession()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	out, err := session.CombinedOutput(command)
	return string(out), err
}

func (this SshClient) Close() {
	for _, client := range this.clients {
		if client == nil {
			continue
		}
		client.Close()
	}
}

func sshConnect(host string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(conf.SSHPrivateKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), sshConfig)
	if err != nil {
		fmt.Println("First ssh attempt failed: " + err.Error())
	}

	return client, err
}
