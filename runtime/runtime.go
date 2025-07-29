package runtime

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"encoding/json"
)

func StartPortProxy(hostPort, containerPort string) {
	ln, err := net.Listen("tcp", ":"+hostPort)
	if err != nil {
		fmt.Println("Failed to listen on port", hostPort, err)
		return
	}
	// fmt.Printf("🔗 Port proxy running: host:%s -> container:%s\n", hostPort, containerPort)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go func(c net.Conn) {
				defer c.Close()
				backend, err := net.Dial("tcp", "127.0.0.1:"+containerPort)
				if err != nil {
					fmt.Println("Failed to connect to container port", containerPort, err)
					return
				}
				defer backend.Close()
				go io.Copy(backend, c)
				io.Copy(c, backend)
			}(conn)
		}
	}()
}

func RunContainer(image, name, port string) error {
	imageDir := filepath.Join("/tmp/underdog/images", image)

	defaultCmd := []string{"/bin/sh"}

	cfg := filepath.Join(imageDir, "config.txt")
	if data, err := os.ReadFile(cfg); err == nil {
		var args []string
		if json.Unmarshal(data, &args) == nil {
			defaultCmd = args
		} else {
			// fallback: parse as plain string
			defaultCmd = strings.Fields(string(data))
		}
	}

	var hostPort, containerPort string
	if port != "" && strings.Contains(port, ":") {
		parts := strings.Split(port, ":")
		hostPort, containerPort = parts[0], parts[1]
		// fmt.Printf("⚠️  Port mapping requested (host:%s -> container:%s)\n", hostPort, containerPort)
		go StartPortProxy(hostPort, containerPort)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// fmt.Println("⚠️  Running in non-isolated mode (Windows)")
		cmd = exec.Command(defaultCmd[0], defaultCmd[1:]...)
		cmd.Dir = imageDir
	} else {
		args := []string{
			"--mount", "--uts", "--ipc", "--net", "--pid", "--fork", "--user", "--map-root-user",
			"chroot", imageDir, defaultCmd[0],
		}
		if len(defaultCmd) > 1 {
			args = append(args, defaultCmd[1:]...)
		}
		cmd = exec.Command("unshare", args...)
	}


	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	return cmd.Run()
}
