package ss

import (
	"errors"
	"fmt"
	"ignite/models"
	"ignite/utils"
	"log"
	"net"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	ImageUrl  string
	client    *docker.Client
	PortRange []int
	Host      string
)

func init() {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Println("New docker client error:", err.Error())
		// TODO panic?
	}
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

func RemoveContaienr(id string) error {
	opt := docker.RemoveContainerOptions{ID: id, RemoveVolumes: true, Force: true}
	return client.RemoveContainer(opt)
}

func StatsOutNet(id string) (uint64, error) {
	errC := make(chan error, 1)
	statsC := make(chan *docker.Stats)
	done := make(chan bool)
	defer close(done)
	go func() {
		errC <- client.Stats(docker.StatsOptions{ID: id, Stats: statsC, Stream: false, Done: done})
		close(errC)
	}()
	stats, ok := <-statsC
	if !ok {
		return 0, errors.New("Can't get stats result")
	}
	err := <-errC
	if err != nil {
		return 0, err
	}
	return stats.Networks["eth0"].TxBytes, nil
}

func CreateAndStart(name string) (*models.ServiceResult, error) {
	r, err := create(name)
	if err != nil {
		return nil, err
	}
	return r, start(r.ID)
}

func getAvaliablePort() (int, error) {
	for port := PortRange[0]; port <= PortRange[1]; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return port, nil
		}
		conn.Close()
	}
	return 0, errors.New("no port available")
}
