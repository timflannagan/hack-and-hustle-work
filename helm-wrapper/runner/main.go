package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	outputBaseDir          = "manifests"
	overlayPlatform        = "openshift"
	meteringFilenamePrefix = "metering-operator"
)

func generateKustomizeManifest(filename string) error {
	cmd := exec.Cmd{
		Path:   "./generate-kustomize-manifest.sh",
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env: []string{
			"KUSTOMIZE_OVERLAY_PATH=../kustomize/platforms/openshift",
			"OUTPUT_PATH=" + filename,
		},
	}
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	f, err := ioutil.TempFile("/tmp/", "manifest")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Remove(f.Name())
		f.Close()
	}()

	err = generateKustomizeManifest(f.Name())
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}

	platformOutputDir := filepath.Join(outputBaseDir, overlayPlatform)
	err = os.MkdirAll(platformOutputDir, 0766)
	if err != nil {
		panic(err)
	}

	// index is the current file #
	// s is the current file buffer
	for _, s := range strings.Split(string(b), "---") {
		var filename string

		for _, line := range strings.Split(s, "\n") {
			if len(line) < 5 {
				continue
			}
			if line[0:5] != "kind:" {
				continue
			}
			filename = fmt.Sprintf("%s-%s.yaml", meteringFilenamePrefix, strings.ToLower(line[6:]))
		}

		err = ioutil.WriteFile(filepath.Join(platformOutputDir, filename), []byte(s), 0755)
		if err != nil {
			panic(err)
		}
	}
}
