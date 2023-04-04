package performance

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func OSvDriver(dockerClient *client.Client, buildContext io.ReadCloser, buildOptions types.ImageBuildOptions) (*PerformanceBenchmark, error) {
	buildOptions.Tags = []string{"osv"}
	buildOptions.Dockerfile = "unikernels/osv.Dockerfile"

	builtImage, err := dockerClient.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(builtImage.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	resp, err := dockerClient.ContainerCreate(context.Background(),
		&container.Config{
			Image: "osv",
			ExposedPorts: nat.PortSet{
				"25565/tcp": struct{}{},
			},
			// AttachStdout: true,
			// AttachStderr: true,
			// Tty:          true,
		},
		&container.HostConfig{
			// Resources: container.Resources{
			// 	Devices: []container.DeviceMapping{
			// 		{PathOnHost: "/dev/kvm", PathInContainer: "/dev/kvm", CgroupPermissions: "rwm"},
			// 	},
			// },
			PortBindings: map[nat.Port][]nat.PortBinding{
				"25565/tcp": {
					{
						HostIP:   "0.0.0.0",
						HostPort: "25565",
					},
				},
			},
			Privileged: true,
		},
		nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	// waiter, _ := dockerClient.ContainerAttach(context.Background(), resp.ID, types.ContainerAttachOptions{
	// 	Stderr: true,
	// 	Stdout: true,
	// 	Stdin:  true,
	// 	Stream: true,
	// })
	// go io.Copy(os.Stdout, waiter.Reader)

	conn, err := net.DialTimeout("tcp", "localhost:25565", 10*time.Second)
	if err != nil {
		return nil, err
	}
	conn.Close()
	boot_start := time.Now()
	println("booting...")

	buf := make([]byte, 1024)
	for string(buf)[0:6] != "booted" {
		time.Sleep(1 * time.Millisecond)
		conn, err = net.DialTimeout("tcp", "localhost:25565", 10*time.Second)
		if err != nil {
			return nil, err
		}
		conn.Read(buf)
	}
	boot_end := time.Now()
	conn.Close()
	println(string(buf))

	statusCh, errCh := dockerClient.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	end := time.Now()

	dockerClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	return &PerformanceBenchmark{
		TimeToBootMs:   boot_end.Sub(boot_start).Milliseconds(),
		TimeToRunMs:    end.Sub(boot_start).Milliseconds(),
		MemoryUsageMiB: 0,
	}, nil
}

func UnikraftDriver(dockerClient *client.Client, buildContext io.ReadCloser, buildOptions types.ImageBuildOptions) (*PerformanceBenchmark, error) {
	return nil, nil
}
