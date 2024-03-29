package docker

import (
	"fmt"

	"github.com/cm-mayfly/cm-mayfly/src/common"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull images of Cloud-Migrator System containers",
	Long:  `Pull images of Cloud-Migrator System containers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Pull images of Cloud-Migrator System containers]")
		fmt.Println()

		if common.DockerFilePath == "" {
			fmt.Println("file is required")
		} else {
			cmdStr := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose -f %s pull", common.CMComposeProjectName, common.DockerFilePath)
			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)
		}

	},
}

func init() {
	dockerCmd.AddCommand(pullCmd)

	pf := pullCmd.PersistentFlags()
	pf.StringVarP(&common.DockerFilePath, "file", "f", common.DefaultDockerComposeConfig, "User-defined configuration file")
	//	cobra.MarkFlagRequired(pf, "file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
