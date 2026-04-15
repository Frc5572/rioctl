package sshclient

import (
	"fmt"
	"regexp"
	"rioctl/internal/utils"
	"slices"
	"strconv"
	"strings"

	"github.com/melbahja/goph"
)

var fileRE = regexp.MustCompile(`^akit_.*_(\w+)_([peq])\d+\.wpilog$`)

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

func (c *Client) ListFiles(dir string) ([]utils.File, error) {
	out, err := c.Run(fmt.Sprintf("ls -l %s | grep .wpilog | awk '{ print $5\",\"$9 }'", dir))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")

	var files []utils.File
	for _, l := range lines {
		if l == "" {
			continue
		}
		eventCode := ""
		matchType := ""
		stringArr := strings.Split(l, ",")
		name := stringArr[1]

		size, err := strconv.ParseInt(strings.TrimSpace(stringArr[0]), 10, 64)
		if err != nil {
			continue
		}
		if fileRE.MatchString(name) {
			matches := fileRE.FindStringSubmatch(name)
			eventCode = matches[1]
			matchType = matches[2]
		}

		files = append(files, utils.File{Name: name, Size: size, Event: eventCode, MatchType: matchType})
	}
	slices.Reverse(files)
	return files, nil
}

func (c *Client) FileSize(path string) (int64, error) {
	out, err := c.Run(fmt.Sprintf("stat -c %%s %s", path))
	if err != nil {
		return 0, err
	}
	size, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
