package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	meteringChartBasePath = "charts/metering-ansible-operator"
	olmDeployBaseDir      = "olm"
	manualDeployBaseDir   = "operator"
	crdBaseDir            = "crds"

	hiveTableFile        = "hive.crd.yaml"
	prestoTableFile      = "prestotable.crd.yaml"
	meteringConfigFile   = "meteringconfig.crd.yaml"
	reportFile           = "report.crd.yaml"
	reportDataSourceFile = "reportdatasource.crd.yaml"
	reportQueryFile      = "reportquery.crd.yaml"
	storageLocationFile  = "storagelocation.crd.yaml"

	// templatesPathPrefix is the base path to
	// the charts/metering-ansible-operator helm
	// chart directory
	templatesPathPrefix    = "templates"
	olmDeployPathPrefix    = "olm"
	manualDeployPathPrefix = "operator"
	crdsPathPrefix         = "crds"

	baseCrdOutputDir = "crds"
)

// Resource is a type responsible for holding the metadata
// required to template a single helm chart into a rendered
// YAML manifest.
type Resource struct {
	Name       string
	InputPath  string
	OutputPath string
	// Optional? bool
}

var errCouldNotFindTemplate error = errors.New("could not find template")

// templateHelmChart is responsible for templating a single helm chart
// Note: for the helm templating + security issues, just mark your $KUBECONFIG
// with only owner can read, write, execute (i.e. chmod 0700 $KUBECONFIG)
func templateHelmChart(values, inputPath, outputPath string) error {
	bin := os.Getenv("HELM_BIN")
	if bin == "" {
		bin = "helm"
	}

	// validate that the `helm` binary in $PATH and if exists,
	// get the full path to that binary
	b, err := exec.LookPath("hack/template-resource.sh")
	if err != nil {
		return err
	}
	d, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine the working dir: %v", err)
	}

	cmd := &exec.Cmd{
		Path:   b,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env: []string{
			"ROOT_DIR=" + d,
			"HELM_BIN=" + bin,
			"CHART=" + meteringChartBasePath,
			"OUTPUT_DIR=" + outputPath,
			"TEMPLATE_RESOURCE_NAME=" + inputPath,
			"DEBUG=false",
			"VALUES_ARGS=-f " + values,
		},
	}
	err = cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), errCouldNotFindTemplate.Error()) {
			return fmt.Errorf("Failed to template helm chart: %+v", err)
		}
		fmt.Println(err)
	}
	f, err := os.Stat(outputPath)
	if err != nil {
		return err
	}

	if f.Size() == 0 {
		err = os.Remove(outputPath)
		if err != nil {
			return err
		}
		fmt.Printf("Skipped generating a YAML manifest for the %s helm chart\n", outputPath)
	}

	return nil
}

func main() {
	baseOutputDir := "manifests"
	err := os.MkdirAll(baseOutputDir, 0755)
	if err != nil {
		panic(err)
	}

	platforms := map[string]string{
		"Openshift": "charts/metering-ansible-operator/values.yaml",
		"Upstream":  "charts/metering-ansible-operator/upstream-values.yaml",
	}
	for platform, values := range platforms {
		platformOutputDir := filepath.Join(baseOutputDir, strings.ToLower(platform))
		err = os.MkdirAll(platformOutputDir, 0755)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Generating %s Manifests\n", platform)
		// TODO: pass a values.yaml chart
		err = templateOLMManifests(values, "4.7.1", platformOutputDir)
		if err != nil {
			panic(err)
		}
		err = templateCrdManifests(values, platformOutputDir)
		if err != nil {
			panic(err)
		}
		err = templateManualManifests(values, platformOutputDir)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Done")
}

func templateOLMManifests(values, csvVersion, outputDir string) error {
	olmOutputDir := filepath.Join(outputDir, "olm")
	err := os.MkdirAll(olmOutputDir, 0755)
	if err != nil {
		return err
	}

	resources := []Resource{
		{
			Name:       "csv",
			InputPath:  "templates/olm/clusterserviceversion.yaml",
			OutputPath: filepath.Join(olmOutputDir, fmt.Sprintf("clusterserviceversion.%s.yaml", csvVersion)),
		},
		{
			Name:       "art",
			InputPath:  "templates/olm/art.yaml",
			OutputPath: filepath.Join(olmOutputDir, "art.yaml"),
		},
		{
			Name:       "catalogsource",
			InputPath:  "templates/olm/catalogsource.yaml",
			OutputPath: filepath.Join(olmOutputDir, "catalogsource.yaml"),
		},
		{
			Name:       "image-references",
			InputPath:  "templates/olm/image-references",
			OutputPath: filepath.Join(olmOutputDir, "image-references"),
		},
		{
			Name:       "operatorgroup",
			InputPath:  "templates/olm/operatorgroup.yaml",
			OutputPath: filepath.Join(olmOutputDir, "operatorgroup.yaml"),
		},
		{
			Name:       "package.yaml",
			InputPath:  "templates/olm/package.yaml",
			OutputPath: filepath.Join(olmOutputDir, "package.yaml"),
		},
		{
			Name:       "subscription",
			InputPath:  "templates/olm/subscription.yaml",
			OutputPath: filepath.Join(olmOutputDir, "subscription.yaml"),
		},
	}

	for _, resource := range resources {
		err = templateHelmChart(values, resource.InputPath, resource.OutputPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func templateManualManifests(values, outputDir string) error {
	operatorOutputDir := filepath.Join(outputDir, "operator")
	err := os.MkdirAll(operatorOutputDir, 0755)
	if err != nil {
		return err
	}

	resources := []Resource{
		{
			Name:       "deployment",
			InputPath:  "templates/operator/deployment.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "deployment.yaml"),
		},
		{
			Name:       "service-account",
			InputPath:  "templates/operator/service-account.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "service-account.yaml"),
		},
		{
			Name:       "meteringconfig",
			InputPath:  "templates/operator/meteringconfig.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "meteringconfig.yaml"),
		},
		{
			Name:       "rolebinding",
			InputPath:  "templates/operator/rolebinding.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "rolebinding.yaml"),
		},
		{
			Name:       "role",
			InputPath:  "templates/operator/role.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "role.yaml"),
		},
		{
			Name:       "clusterrolebinding",
			InputPath:  "templates/operator/clusterrolebinding.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "clusterrolebinding.yaml"),
		},
		{
			Name:       "clusterrole",
			InputPath:  "templates/operator/clusterrole.yaml",
			OutputPath: filepath.Join(operatorOutputDir, "clusterrole.yaml"),
		},
	}

	for _, resource := range resources {
		err = templateHelmChart(values, resource.InputPath, resource.OutputPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func templateCrdManifests(values, baseOutputDir string) error {
	crdBaseDir := filepath.Join(templatesPathPrefix, crdsPathPrefix)

	crdOutputDir := filepath.Join(baseOutputDir, "crds")
	err := os.MkdirAll(crdOutputDir, 0755)
	if err != nil {
		return err
	}

	// pathToCrdResources is an array of the metering-related CRD
	// resources where the string key is the metadata.name of the CRD
	// and the key is the corresponding output directory once that YAML
	// manifest has been templated using the $HELM_BIN helm binary.
	crds := []Resource{
		{
			Name:       "HiveTable",
			InputPath:  filepath.Join(crdBaseDir, hiveTableFile),
			OutputPath: filepath.Join(crdOutputDir, hiveTableFile),
		},
		{
			Name:       "PrestoTable",
			InputPath:  filepath.Join(crdBaseDir, prestoTableFile),
			OutputPath: filepath.Join(crdOutputDir, prestoTableFile),
		},
		{
			Name:       "MeteringConfig",
			InputPath:  filepath.Join(crdBaseDir, meteringConfigFile),
			OutputPath: filepath.Join(crdOutputDir, meteringConfigFile),
		},
		{
			Name:       "Report",
			InputPath:  filepath.Join(crdBaseDir, reportFile),
			OutputPath: filepath.Join(crdOutputDir, reportFile),
		},
		{
			Name:       "ReportDataSource",
			InputPath:  filepath.Join(crdBaseDir, reportDataSourceFile),
			OutputPath: filepath.Join(crdOutputDir, reportDataSourceFile),
		},
		{
			Name:       "ReportQuery",
			InputPath:  filepath.Join(crdBaseDir, reportQueryFile),
			OutputPath: filepath.Join(crdOutputDir, reportQueryFile),
		},
		{
			Name:       "StorageLocation",
			InputPath:  filepath.Join(crdBaseDir, storageLocationFile),
			OutputPath: filepath.Join(crdOutputDir, storageLocationFile),
		},
	}

	// TODO: make parallel
	// TODO: should this be an array of errors?
	for _, crd := range crds {
		err := templateHelmChart(values, crd.InputPath, crd.OutputPath)
		if err != nil {
			return err
		}
	}
	return nil
}
