package main

import (
	"errors"
	"fmt"
	_ "io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Note: /usr/bin/diff returns an exit status of 0 if the
// two files are equal, 1 if not, and 2 if something unexpected
// occurred. In our case, we don't care about anything besides
// a return code of 2.
var errDiffTroublingErrorCode = errors.New("exit status 2")

func main() {
	op, err := buildResources("../manifests/openshift/operator/")
	if err != nil {
		panic(err)
	}

	up, err := buildResources("../manifests/upstream/operator/")
	if err != nil {
		panic(err)
	}

	err = compareManifestArrays("../manifests", op, up)
	if err != nil {
		panic(err)
	}
}

func buildResources(manifestPath string) ([]string, error) {
	var resources []string

	// walk through the @manifestPath directory, appending any non-directory
	// files to the resources array of filenames
	err := filepath.Walk(manifestPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		absPath, err := filepath.Abs(filepath.Join(manifestPath, f.Name()))
		if err != nil {
			return err
		}

		resources = append(resources, absPath)
		return nil
	})
	if err != nil {
		return resources, err
	}

	return resources, nil
}

func compareManifestArrays(basePath string, arr1 []string, arr2 []string) error {
	// TODO: eventually stop enforcing this check
	if len(arr1) != len(arr2) {
		return fmt.Errorf("The length of the array parameters need to be equal")
	}

	// TODO: this is such a poor implementation
	for i, elem := range arr1 {
		fmt.Printf("Comparing %s with %s\n", elem, arr2[i])

		cmd := exec.Command("/usr/bin/diff", elem, arr2[i])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil && errors.Is(err, errDiffTroublingErrorCode) {
			return err
		}
	}

	return nil
}
