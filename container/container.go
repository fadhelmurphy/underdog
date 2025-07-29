package container

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"underdog/parser"
)

func tarDirectory(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tw := tar.NewWriter(tarfile)
	defer tw.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == source {
			return nil
		}
		relPath := strings.TrimPrefix(path, source)
		if relPath == "" {
			return nil
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			io.Copy(tw, file)
		}
		return nil
	})
	return nil
}

func BuildImage(image string, stages [][]parser.Instruction) error {
	imageDir := filepath.Join("/tmp/underdog/images", image)
	rootfsDir := filepath.Join(imageDir, "rootfs")
	os.MkdirAll(rootfsDir, 0755)
	workingDir := rootfsDir

	var defaultCmd []string
	var workDir = "/"

	for _, stage := range stages {
		for _, inst := range stage {
			switch inst.Cmd {
			case "COPY":
				src := inst.Args[0]
				dest := inst.Args[1]
				if dest == "." || strings.HasSuffix(dest, "/") {
					dest = filepath.Join(workingDir, filepath.Base(src))
				} else {
					dest = filepath.Join(workingDir, dest)
				}
				data, err := os.ReadFile(src)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", src, err)
				}
				os.MkdirAll(filepath.Dir(dest), 0755)
				err = os.WriteFile(dest, data, 0644)
				if err != nil {
					return fmt.Errorf("failed to write file %s: %w", dest, err)
				}
			case "RUN":
				cmd := exec.Command("sh", "-c", strings.Join(inst.Args, " "))
				cmd.Dir = workingDir
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("run command failed: %w", err)
				}
			case "WORKDIR":
				workDir = inst.Args[0]
				workingDir = filepath.Join(rootfsDir, inst.Args[0])
				os.MkdirAll(workingDir, 0755)
			case "CMD", "ENTRYPOINT":
				defaultCmd = inst.Args
			}
		}
	}

	// Tarball layer
	blobsDir := filepath.Join(imageDir, "blobs")
	os.MkdirAll(blobsDir, 0755)
	tarPath := filepath.Join(blobsDir, "layer.tar")
	if err := tarDirectory(rootfsDir, tarPath); err != nil {
		return fmt.Errorf("failed to create layer tar: %w", err)
	}

	// Config JSON
	configContent := fmt.Sprintf(`{
  "architecture": "amd64",
  "os": "linux",
  "config": {
    "Cmd": ["%s"],
    "Entrypoint": null,
    "WorkingDir": "%s"
  }
}`, strings.Join(defaultCmd, `","`), workDir)
	os.WriteFile(filepath.Join(imageDir, "config.json"), []byte(configContent), 0644)

	// Manifest JSON
	manifestContent := `{
  "schemaVersion": 2,
  "config": "config.json",
  "layers": ["blobs/layer.tar"]
}`
	os.WriteFile(filepath.Join(imageDir, "manifest.json"), []byte(manifestContent), 0644)

	fmt.Println("Image built at", imageDir)
	return nil
}
