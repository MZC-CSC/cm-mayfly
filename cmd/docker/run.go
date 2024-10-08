package docker

import (
	"fmt"

	"github.com/cm-mayfly/cm-mayfly/common"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Setup and Run Cloud-Migrator System",
	Long:  `Setup and Run Cloud-Migrator System`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Setup and Run Cloud-Migrator]")
		fmt.Println()

		if DockerFilePath == "" {
			fmt.Println("--file (-f) argument is required but not provided.")
		} else {
			//cmdStr := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose -f %s up -d", CMComposeProjectName, DockerFilePath)
			cmdStr := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose -f %s up ", CMComposeProjectName, DockerFilePath)

			// If there are additional arguments, treat them as services or additional commands and add them to the existing command with an additional
			if len(args) > 0 {
				cmdStr += args[0]

				// Explicitly passing the service name as a filter (--service) option or argument would be fine.
				// serviceName := args[0]
				// cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose -f %s up %s", CMComposeProjectName, DockerFilePath, serviceName)
			}

			// fmt.Println(cmdStr)
			common.SysCall(cmdStr)
		}
	},
}

func init() {
	dockerCmd.AddCommand(runCmd)

	pf := runCmd.PersistentFlags()
	pf.StringVarP(&DockerFilePath, "file", "f", DefaultDockerComposeConfig, "User-defined configuration file")
	//pf.StringVarP(&k8sprovider, "k8sprovider", "", common.NotDefined, "Kind of Managed K8s services") //@todo

	// runCmd.MarkPersistentFlagRequired("k8sprovider")

	//	cobra.MarkFlagRequired(pf, "file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
