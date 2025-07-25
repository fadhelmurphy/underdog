package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"underdog/parser"
	"strings"
)

func BuildImage(image string, stages [][]parser.Instruction) error {
	imageDir := filepath.Join("/tmp/underdog/images", image)
	os.MkdirAll(imageDir, 0755)

	workingDir := imageDir

	for _, stage := range stages {
		for _, inst := range stage {
			switch inst.Cmd {
			case "COPY":
				src := inst.Args[0]
				dest := inst.Args[1]

				if dest == "." || strings.HasSuffix(dest, "/") {
					dest = filepath.Join(imageDir, filepath.Base(src))
				} else {
					dest = filepath.Join(imageDir, dest)
				}

				data, err := os.ReadFile(src)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", src, err)
				}
				err = os.WriteFile(dest, data, 0644)
				if err != nil {
					return fmt.Errorf("failed to write file %s: %w", dest, err)
				}
			case "RUN":
				cmd := exec.Command("sh", "-c", strings.Join(inst.Args, " "))
				cmd.Dir = workingDir
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("run command failed: %w", err)
				}
			case "WORKDIR":
				workingDir = filepath.Join(imageDir, inst.Args[0])
				os.MkdirAll(workingDir, 0755)
			}
		}
	}
	fmt.Println("Image built at", imageDir)
	return nil
}
