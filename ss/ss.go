package ss

import (
	"errors"
	"fmt"
	"ignite/models"
	"ignite/utils"
	"net"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	ImageUrl  string
	client    *docker.Client
	portRange = []int{4000, 6000}
)

func Init() (err error) {
	client, err = docker.NewClientFromEnv()
	return err
}

func create(name string) (*models.ServiceResult, error) {
	password := utils.NewPasswd(16)
	port, err := getAvaliablePort()
	if err != nil {
		return nil, err
	}
	portStr := fmt.Sprintf("%d", port)
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image:        ImageUrl,
			Cmd:          []string{"-k", password, "-p", portStr},
			ExposedPorts: map[docker.Port]struct{}{docker.Port(portStr + "/tcp"): {}},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(portStr + "/tcp"): []docker.PortBinding{{HostPort: portStr}}},
			RestartPolicy: docker.AlwaysRestart(),
		},
	})
	if err != nil {
		return nil, err
	}
	r := &models.ServiceResult{
		ID:       container.ID,
		Password: password,
		Port:     port,
	}
	return r, nil
}

func start(id string) error {
	return client.StartContainer(id, nil)
}

func CreateAndStart(name string) (*models.ServiceResult, error) {
	r, err := create(name)
	if err != nil {
		return nil, err
	}
	return r, start(r.ID)
}

func getAvaliablePort() (int, error) {
	for port := portRange[0]; port <= portRange[1]; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return port, nil
		}
		conn.Close()
	}
	return 0, errors.New("no port available")
}
