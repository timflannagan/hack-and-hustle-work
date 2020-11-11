package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	platform  string
	outputDir string

	rootCmd = &cobra.Command{
		Use:   "generate-manifests",
		Short: "Generate the Metering deployment manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Generate the Metering deployment manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerateManifests()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&platform, "platform", "", "The Kubernetes platform to generate YAML manifests for")
	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "manifests/deploy/", "The base platform to store rendered YAML manifests in")
}

func main() {
	rootCmd.AddCommand(runCmd)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		panic(err)
	}
}

func runGenerateManifests() error {
	fmt.Println("this ran")
	return nil
}
