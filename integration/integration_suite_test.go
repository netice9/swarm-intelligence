package integration_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var docker *client.Client

func startContainer(containerName, imageName string, cmd []string, links []string, portMapping map[int]int) {

	_, _, err := docker.ImageInspectWithRaw(context.Background(), imageName)

	if client.IsErrNotFound(err) {
		var output io.ReadCloser
		output, err = docker.ImagePull(context.Background(),
			imageName,
			types.ImagePullOptions{},
		)
		Expect(err).ToNot(HaveOccurred())
		io.Copy(os.Stdout, output)
		output.Close()
	}

	Expect(err).ToNot(HaveOccurred())

	cont, err := docker.ContainerInspect(context.Background(), containerName)

	if err == nil {
		err = docker.ContainerRemove(
			context.Background(),
			cont.ID,
			types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			},
		)
		Expect(err).ToNot(HaveOccurred())
	}

	portMap := nat.PortMap{}

	exposed := map[nat.Port]struct{}{}

	for h, c := range portMapping {
		portMap[nat.Port(fmt.Sprintf("%d/tcp", c))] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", h)}}
		exposed[nat.Port(fmt.Sprintf("%d/tcp", c))] = struct{}{}
	}

	created, err := docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        imageName,
			Env:          []string{},
			Cmd:          cmd,
			ExposedPorts: exposed,
		},
		&container.HostConfig{
			Privileged:   true,
			Links:        links,
			PortBindings: portMap,
		},
		&network.NetworkingConfig{},
		containerName,
	)
	Expect(err).ToNot(HaveOccurred())

	err = docker.ContainerStart(
		context.Background(),
		created.ID,
		types.ContainerStartOptions{},
	)
	Expect(err).ToNot(HaveOccurred())
}

var _ = BeforeSuite(func() {
	var err error
	docker, err = client.NewEnvClient()
	Expect(err).ToNot(HaveOccurred())

	startContainer(
		"test-registry",
		"registry:2.5.1",
		[]string{},
		[]string{},
		map[int]int{5000: 5000},
	)

	startContainer(
		"swarm1",
		"docker:1.12.3",
		[]string{
			"docker",
			"daemon",
			"--insecure-registry=test-registry",
			"-H", "tcp://0.0.0.0:2375",
			"-H", "unix:///var/run/docker.sock",
			"--iptables=false",
		},
		[]string{"test-registry"},
		map[int]int{5001: 2375},
	)

	startContainer(
		"swarm2",
		"docker:1.12.3",
		[]string{
			"docker",
			"daemon",
			"--insecure-registry=test-registry",
			"-H", "tcp://0.0.0.0:2375",
			"-H", "unix:///var/run/docker.sock",
			"--iptables=false",
		},
		[]string{"test-registry", "swarm1"},
		map[int]int{5002: 2375},
	)

	sw1Docker, err := client.NewClient("tcp://localhost:5001", client.DefaultVersion, nil, nil)
	Expect(err).ToNot(HaveOccurred())
	nodeID, err := sw1Docker.SwarmInit(context.Background(), swarm.InitRequest{
		ListenAddr: "0.0.0.0:2377",
	})
	Expect(err).ToNot(HaveOccurred())
	log.Println("nodeID", nodeID)

	sw, err := sw1Docker.SwarmInspect(context.Background())
	Expect(err).ToNot(HaveOccurred())

	sw2Docker, err := client.NewClient("tcp://localhost:5002", client.DefaultVersion, nil, nil)
	Expect(err).ToNot(HaveOccurred())
	err = sw2Docker.SwarmJoin(context.Background(), swarm.JoinRequest{
		JoinToken:   sw.JoinTokens.Manager,
		ListenAddr:  "0.0.0.0:2377",
		RemoteAddrs: []string{"swarm1:2377"},
	})
	Expect(err).ToNot(HaveOccurred())

})
