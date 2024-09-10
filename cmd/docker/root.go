/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package docker

import (
	"fmt"

	"github.com/cm-mayfly/cm-mayfly/cmd"
	"github.com/spf13/cobra"
)

// restCmd represents the rest command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Installing and managing cloud-migrator's infrastructure",
	Long: `Build the environment of the infrastructure required for cloud-migrator and monitor the running status of the infrastructure.
For example, you can setup and run, stop, and ... Cloud-Migrator runtimes.

- ./mayfly docker pull [-f ./conf/docker/docker-compose.yaml]
- ./mayfly docker run [-f ./conf/docker/docker-compose.yaml]
- ./mayfly docker info
- ./mayfly docker stop [-f ./conf/docker/docker-compose.yaml]
- ./mayfly docker remove [-f ./conf/docker/docker-compose.yaml] -v -i

	     `,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(cmd.UsageString())
		fmt.Println(cmd.Help())
	},
}

func init() {
	cmd.RootCmd.AddCommand(dockerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
