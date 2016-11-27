package integration_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"testing"
)

var agoutiDriver *agouti.WebDriver

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

func createSwarm() {
	startContainer(
		"test-registry",
		"registry:2.5.1",
		[]string{},
		[]string{},
		map[int]int{5000: 5000},
	)

	startContainer(
		"swarm1",
		"docker:1.12.3-dind",
		[]string{
			"docker",
			"daemon",
			"--insecure-registry=test-registry:5000",
			"-H", "tcp://0.0.0.0:2375",
			"-H", "unix:///var/run/docker.sock",
			// "--iptables=false",
		},
		[]string{"test-registry"},
		map[int]int{5001: 2375, 5080: 3000},
	)

	startContainer(
		"swarm2",
		"docker:1.12.3-dind",
		[]string{
			"docker",
			"daemon",
			"--insecure-registry=test-registry:5000",
			"-H", "tcp://0.0.0.0:2375",
			"-H", "unix:///var/run/docker.sock",
			// "--iptables=false",
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

}

func pushImage() {
	wd, err := os.Getwd()
	Expect(err).ToNot(HaveOccurred())
	dir := filepath.Dir(wd)

	ctxReader, ctxWriter := io.Pipe()

	go func() {
		Expect(tarit(dir, ctxWriter)).To(Succeed())
		Expect(ctxWriter.Close()).To(Succeed())
	}()

	b, err := docker.ImageBuild(
		context.Background(),
		ctxReader,
		types.ImageBuildOptions{
			Tags:       []string{"localhost:5000/si/swarm-intelligence:current"},
			Dockerfile: "Dockerfile.integration",
			Squash:     true,
		},
	)
	Expect(err).ToNot(HaveOccurred())
	_, err = io.Copy(os.Stdout, b.Body)
	Expect(err).ToNot(HaveOccurred())
	Expect(b.Body.Close()).To(Succeed())

	out, err := docker.ImagePush(context.Background(), "localhost:5000/si/swarm-intelligence:current", types.ImagePushOptions{
		RegistryAuth: "-",
	})
	Expect(err).ToNot(HaveOccurred())
	_, err = io.Copy(os.Stdout, out)
	Expect(err).ToNot(HaveOccurred())
	Expect(out.Close()).To(Succeed())
}

func uint64Ptr(val uint64) *uint64 {
	return &val
}

func deploySwarmIntelligenceService() {
	sw1Docker, err := client.NewClient("tcp://localhost:5001", client.DefaultVersion, nil, nil)
	Expect(err).ToNot(HaveOccurred())

	resp, err := sw1Docker.ServiceCreate(

		context.Background(),
		swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "swarm-intelligence-head",
			},
			Mode: swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{
					Replicas: uint64Ptr(1),
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: "test-registry:5000/si/swarm-intelligence:current",
					Env: []string{
						"PORT=3000",
					},
					Mounts: []mount.Mount{{
						Type:     "bind",
						Source:   "/var/run/docker.sock",
						Target:   "/var/run/docker.sock",
						ReadOnly: false,
					},
					},
				},
				Resources: &swarm.ResourceRequirements{},
				// Placement: &swarm.Placement{},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: []swarm.PortConfig{
					{
						Name:          "web",
						Protocol:      swarm.PortConfigProtocolTCP,
						TargetPort:    3000,
						PublishedPort: 3000,
					},
				},
				Mode: swarm.ResolutionModeVIP,
			},
		},
		types.ServiceCreateOptions{},
	)
	Expect(err).ToNot(HaveOccurred())
	log.Println("serviceID", resp.ID)
}

var _ = BeforeSuite(func() {
	var err error
	docker, err = client.NewEnvClient()
	Expect(err).ToNot(HaveOccurred())
	createSwarm()
	pushImage()
	deploySwarmIntelligenceService()

	agoutiDriver = agouti.PhantomJS()
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
})
