package sshclient

import (
	"fmt"
	"strings"

	"github.com/melbahja/goph"
)

type Client struct {
	client *goph.Client
}

func (c *Client) Raw() *goph.Client {
	return c.client
}

func New(host, user, password string) (*Client, error) {
	auth := goph.Password(password)

	client, err := goph.NewUnknown(user, host, auth)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Run(cmd string) (string, error) {
	out, err := c.client.Run(cmd)
	return string(out), err
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) ListFiles(dir string) ([]string, error) {
	out, err := c.Run(fmt.Sprintf("find %s -type f -name '*.wpilog'", dir))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")

	var files []string
	for _, l := range lines {
		if l == "" {
			continue
		}
		files = append(files, strings.TrimPrefix(l, dir+"/"))
	}

	return files, nil
}
