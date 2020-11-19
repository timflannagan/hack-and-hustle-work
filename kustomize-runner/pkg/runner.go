package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/cli-runtime/pkg/kustomize"
	"sigs.k8s.io/kustomize/pkg/fs"
)

const (
	kustomizeBaseDir = "kustomize"
	overlayBaseDir   = "platforms"

	outputBaseDir          = "manifests"
	meteringFilenamePrefix = "metering-operator"
)

type Runner struct {
	Platforms         []string
	KustomizeBaseDir  string
	ManifestOutputDir string
}

func NewRunner(kustomizeBaseDir, manifestOutputDir string, platforms []string) *Runner {
	return &Runner{
		Platforms:         platforms,
		KustomizeBaseDir:  kustomizeBaseDir,
		ManifestOutputDir: manifestOutputDir,
	}
}

func (r *Runner) Run() error {
	out := bytes.Buffer{}
	filesys := fs.MakeRealFS()

	for _, platform := range r.Platforms {
		fmt.Printf("Processing the %s platform\n", platform)

		// should probably stat this before punting to kustomize
		err := kustomize.RunKustomizeBuild(&out, filesys, filepath.Join(kustomizeBaseDir, overlayBaseDir, platform))
		if err != nil {
			panic(fmt.Errorf("Failed to run kustomize: %+v", err))
		}

		platformOutputDir := filepath.Join(outputBaseDir, platform)
		err = os.MkdirAll(platformOutputDir, 0644)
		if err != nil {
			panic(err)
		}

		err = splitKustomizeManifest(platformOutputDir, out.Bytes())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Finished rendering YAML manifests")
	return nil
}

func splitKustomizeManifest(outputDir string, manifests []byte) error {
	const kindPropertyMatch = "kind:"

	// Split the string representation of the @manifests bytes array
	// by the `---` deliminitor. This results in individual manifests
	// that don't contain the deliminitor we care about such that we
	// can store those rendered manifests on disk.
	for _, manifest := range strings.Split(string(manifests), "---") {
		var filename string

		// Continue processing the string YAML manifests, searching
		// for the top-level `kind: ...` so we can use that within the
		// output filename. If the length of the line is less than five,
		// continue to the next iterator to avoid a zero-indexing error,
		// else, find the first match and exit early.
		for _, line := range strings.Split(manifest, "\n") {
			if len(line) < 5 {
				continue
			}
			if line[0:5] == kindPropertyMatch {
				filename = fmt.Sprintf("%s-%s.yaml", meteringFilenamePrefix, strings.ToLower(line[6:]))
				break
			}
		}

		err := ioutil.WriteFile(filepath.Join(outputDir, filename), []byte(manifest), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write to the %s filename: %v", filename)
		}
	}

	return nil
}
