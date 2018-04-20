package ss

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/utils"
)

const (
	SS_IMAGE  = "goignite/ss-libev:latest"
	SSR_IMAGE = "goignite/ssr:latest"
)

var (
	client    *docker.Client
	PortRange []int
	Host      string
)

func init() {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Fatalln("ignite require Docker installed")
	}
}

func CreateContainer(serverType, name, method, password string, port int) (*models.ServiceResult, error) {
	image := ""
	switch serverType {
	case "SS":
		image = SS_IMAGE
	case "SSR":
		image = SSR_IMAGE
	default:
		return nil, errors.New("invalid server type")
	}
	PullImage(image, true)
	if password == "" {
		password = utils.NewPasswd(16)
	}
	portStr := fmt.Sprintf("%d", port)
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
			Cmd:   []string{"-k", password, "-m", method},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port("3389/tcp"): {{HostPort: portStr}},
				docker.Port("3389/udp"): {{HostPort: portStr}},
			},
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

func StartContainer(id string) error {
	return client.StartContainer(id, &docker.HostConfig{})
}

func PullImage(image string, quiet ...bool) error {
	var output io.Writer
	output = os.Stdout
	if len(quiet) > 0 && quiet[0] {
		output = ioutil.Discard
	}

	return client.PullImage(docker.PullImageOptions{Repository: image, OutputStream: output},
		docker.AuthConfiguration{})
}

func KillContainer(id string) error {
	opt := docker.KillContainerOptions{ID: id}
	return client.KillContainer(opt)
}

func StopContainer(id string, timeout ...uint) error {
	var t uint = 10
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return client.StopContainer(id, t)
}

func RemoveContainer(id string) error {
	opt := docker.RemoveContainerOptions{ID: id, RemoveVolumes: true, Force: true}
	err := client.RemoveContainer(opt)
	if err != nil {
		if _, ok := err.(*docker.NoSuchContainer); ok {
			return nil
		}
	}
	return err
}

func IsContainerRunning(id string) bool {
	info, err := client.InspectContainer(id)
	if err != nil {
		return false
	}

	return info.State.Running
}

func GetContainerStartTime(id string) (*time.Time, error) {
	info, err := client.InspectContainer(id)
	if err != nil {
		return nil, err
	}
	return &info.State.StartedAt, nil
}

func GetContainerStatsOutNet(id string) (uint64, error) {
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

func CreateAndStartContainer(serverType, name, method, password string, port int) (*models.ServiceResult, error) {
	r, err := CreateContainer(serverType, name, method, password, port)
	if err != nil {
		return nil, err
	}
	return r, StartContainer(r.ID)
}

func ContainerExist(id string) (bool, error) {
	_, err := client.InspectContainer(id)
	if err == nil {
		return true, nil
	}
	if _, ok := err.(*docker.NoSuchContainer); ok {
		return false, nil
	}
	return false, err
}
