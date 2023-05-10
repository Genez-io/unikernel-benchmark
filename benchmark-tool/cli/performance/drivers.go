package performance

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

func connectToTCPServer() (*net.Conn, error) {
	var conn net.Conn
	var err error

	for i := 0; ; i++ {
		conn, err = net.DialTimeout("tcp", "localhost"+DOCKER_PORT, 10*time.Second)
		if err != nil {
			if i >= MAX_RETRIES {
				return nil, err
			}

			time.Sleep(5 * time.Millisecond)
			continue
		}
		conn.SetDeadline(time.Time{})

		break
	}

	return &conn, nil
}

func readTCPMessage(conn net.Conn) ([]byte, int, error) {
	var n int
	var err error
	data := make([]byte, 0)

	// Read length of message
	for bytes_recv := 0; bytes_recv < 4; bytes_recv += n {
		buffer := make([]byte, 4-bytes_recv)
		n, err = conn.Read(buffer)
		if err != nil {
			return nil, 0, err
		}

		data = append(data, buffer[:n]...)
	}
	length := int(binary.LittleEndian.Uint32(data))

	data = make([]byte, 0)
	// Read message
	for bytes_recv := 0; bytes_recv < length; bytes_recv += n {
		buffer := make([]byte, length-bytes_recv)
		n, err = conn.Read(buffer)
		if err != nil {
			return nil, 0, err
		}

		data = append(data, buffer[:n]...)
	}

	return data, length, nil
}

func getStaticMetricsAndBoot(conn net.Conn) (*StaticMetrics, error) {
	var staticMetrics StaticMetrics
	var err error
	var buffer []byte
	var bytes_read int

	buffer, bytes_read, err = readTCPMessage(conn)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Static metrics received: " + string(buffer[:bytes_read]))

	err = json.Unmarshal(buffer[:bytes_read], &staticMetrics)
	if err != nil {
		return nil, err
	}

	buffer, bytes_read, err = readTCPMessage(conn)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Start booting message received: " + string(buffer[:bytes_read]))

	if string(buffer[:bytes_read]) != "start_booting" {
		return nil, errors.New("Failed to boot unikernel")
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

func waitUnikernelExecutionEnd(conn net.Conn) error {
	buffer, bytes_read, err := readTCPMessage(conn)
	if err != nil {
		return err
	}
	logrus.Debug("Execution end message received: " + string(buffer[:bytes_read]))

	if string(buffer[:bytes_read]) != "execution_ended" {
		return errors.New("Failed to wait unikernel execution end")
	}

	return nil
}

func getRuntimeMetrics(conn net.Conn) (*RuntimeMetrics, error) {
	var runtimeMetrics RuntimeMetrics

	buffer, bytes_read, err := readTCPMessage(conn)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Runtime metrics received: " + string(buffer[:bytes_read]))

	err = json.Unmarshal(buffer[:bytes_read], &runtimeMetrics)
	if err != nil {
		return nil, errors.New("Failed to unmarshal runtime metrics")
	}

	return &runtimeMetrics, nil
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
		line := scanner.Text()
		logrus.Debug(line)
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
	go func() {
		// Read from waiter.Reader and log it to DEBUG
		scanner := bufio.NewScanner(waiter.Reader)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Replace(line, "\r", "", -1)

			logrus.Debug(line)
		}
	}()

	time.Sleep(time.Second)
	tcpConn, err := connectToTCPServer()

	staticMetrics, err := getStaticMetricsAndBoot(*tcpConn)
	if err != nil {
		return nil, err
	}

	boot_start := time.Now()
	err = waitUnikernetToBoot()
	boot_end := time.Now()

	logrus.Info("Unikernel booted, waiting for execution to end...")

	if err != nil {
		return nil, err
	}

	err = waitUnikernelExecutionEnd(*tcpConn)
	end := time.Now()

	runtimeMetrics, err := getRuntimeMetrics(*tcpConn)
	if err != nil {
		return nil, err
	}

	statusCh, errCh := dockerClient.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	dockerClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	return &PerformanceBenchmark{
		TimeToBootMs:   boot_end.Sub(boot_start).Milliseconds(),
		TimeToRunMs:    end.Sub(boot_start).Milliseconds(),
		StaticMetrics:  *staticMetrics,
		RuntimeMetrics: *runtimeMetrics,
	}, nil
}

func UnikraftDriver(dockerClient *client.Client, buildContext io.ReadCloser, buildOptions types.ImageBuildOptions) (*PerformanceBenchmark, error) {
	return nil, nil
}
