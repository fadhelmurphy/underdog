package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunContainer(image, name, port string) error {
	imageDir := filepath.Join("/tmp/underdog/images", image)
	cmd := exec.Command("unshare", "--mount", "--uts", "--ipc", "--net", "--pid", "--fork", "--user", "--map-root-user", "chroot", imageDir, "/bin/sh", "/entrypoint.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	if port != "" {
		fmt.Println("⚠️  Port mapping not yet implemented in user-space mode.")
	}

	return cmd.Run()
}
