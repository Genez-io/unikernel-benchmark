package performance

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func getStaticMetricsAndBoot() (*StaticMetrics, error) {
	var conn net.Conn
	var err error
	var staticMetrics StaticMetrics

	for i := 0; ; i++ {
		conn, err = net.DialTimeout("tcp", "localhost"+DOCKER_PORT, 10*time.Second)
		if err != nil {
			if i >= MAX_RETRIES {
				return nil, err
			}

			time.Sleep(5 * time.Millisecond)
			continue
		}

		break
	}

	buffer := make([]byte, 1024)
	bytes_read, _ := conn.Read(buffer)

	err = json.Unmarshal(buffer[:bytes_read], &staticMetrics)
	if err != nil {
		return nil, err
	}

	for bytes_read, _ := conn.Read(buffer); bytes_read != 0; {
		time.Sleep(time.Millisecond)
	}

	return &staticMetrics, nil
}

func waitUnikernetToBoot() error {
	var conn net.Conn
	var err error

	udpServer, err := net.ResolveUDPAddr("udp", DOCKER_PORT)
	if err != nil {
		return err
	}

	for i := 0; ; i++ {
		conn, err = net.DialUDP("udp", nil, udpServer)
		if err != nil {
			if i >= MAX_RETRIES {
				return err
			}

			time.Sleep(5 * time.Millisecond)
			continue
		}

		break
	}

	buf := make([]byte, 1024)
	for string(buf)[0:6] != "booted" {
		_, err = conn.Write([]byte("booting..."))
		if err != nil {
			time.Sleep(time.Millisecond)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		conn.Read(buf)
	}
	conn.Close()

	return nil
}

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
				"25565/udp": struct{}{},
			},
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				"25565/tcp": {
					{
						HostIP:   "0.0.0.0",
						HostPort: "25565",
					},
				},
				"25565/udp": {
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

	waiter, _ := dockerClient.ContainerAttach(context.Background(), resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Stream: true,
	})
	go io.Copy(os.Stdout, waiter.Reader)

	time.Sleep(time.Second)
	staticMetrics, err := getStaticMetricsAndBoot()
	if err != nil {
		return nil, err
	}

	// go func() {
	// 	i := 0
	// 	for {
	// 		println(fmt.Sprint(i) + "ms")
	// 		time.Sleep(10 * time.Millisecond)
	// 		i += 10
	// 	}
	// }()

	boot_start := time.Now()
	waitUnikernetToBoot()
	boot_end := time.Now()

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
		StaticMetrics:  *staticMetrics,
	}, nil
}

func UnikraftDriver(dockerClient *client.Client, buildContext io.ReadCloser, buildOptions types.ImageBuildOptions) (*PerformanceBenchmark, error) {
	return nil, nil
}
