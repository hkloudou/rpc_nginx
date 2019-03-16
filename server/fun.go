package main

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

//NewDockerClient NewDockerClient
func NewDockerClient(endpoint string) (*docker.Client, error) {
	if strings.HasPrefix(endpoint, "unix:") {
		return docker.NewClient(endpoint)
	}
	return docker.NewClient(endpoint)
}
