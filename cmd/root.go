package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"

	lib "github.com/justmiles/glue-deploy/lib"
)

var options lib.DeploymentOptions

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glue-deploy",
	Short: "An opinionated deployment process for AWS Glue",
	Run: func(rootCmd *cobra.Command, args []string) {

		options.Sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// create a new deployment
		d := lib.NewDeployment(&options)
		// get deployment changes
		err := d.Build()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// exit if no changes
		if !d.HasChanges() {
			fmt.Println("No Changes")
			os.Exit(0)
		}

		// display to user
		fmt.Println(d)

		// handle approval
		if !options.AutoApprove {
			if !confirmDeployment(d) {
				fmt.Println("Deployment Cancelled")
				os.Exit(130)
			}
		}

		// deploy
		fmt.Print("Deploying..")

		err = d.Deploy()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Print("..Done!\n")
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}}
`)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&options.Artifact, "artifact-id", "a", "", "id artifact to deploy")
	rootCmd.MarkFlagRequired("artifact-id")

	rootCmd.Flags().StringVarP(&options.Version, "artifact-version", "v", "", "artifact version to set")
	rootCmd.MarkFlagRequired("artifact-version")

	rootCmd.Flags().StringVarP(&options.Environment, "environment", "e", "", "target deployment environment")
	rootCmd.MarkFlagRequired("environment")

	rootCmd.Flags().StringVarP(&options.Role, "role", "r", "", "(optional) an IAM role ARN to assume before invoking a deployment")
	rootCmd.Flags().BoolVar(&options.AutoApprove, "auto-approve", false, "automatically approve version changes")

}

func confirmDeployment(d *lib.Deployment) bool {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Approve Changes [y/n]: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(128)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
