package framework

import (
	"fmt"
	"strings"

	root "github.com/cm-mayfly/cm-mayfly/src/cmd"

	"github.com/cm-mayfly/cm-mayfly/src/common"
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

		if common.FileStr == "" {
			fmt.Println("--file (-f) argument is required but not provided.")
		} else {
			common.FileStr = common.GenConfigPath(common.FileStr, common.CMMayflyMode)

			var cmdStr string
			switch common.CMMayflyMode {
			case common.ModeDockerCompose:
				cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose -f %s up", common.CMComposeProjectName, common.FileStr)
				//fmt.Println(cmdStr)
				common.SysCall(cmdStr)
			case common.ModeKubernetes:
				if root.K8sprovider == common.NotDefined {
					fmt.Print(`--k8sprovider argument is required but not provided.
					e.g.
					--k8sprovider=gke
					--k8sprovider=eks
					--k8sprovider=aks
					--k8sprovider=mcks
					--k8sprovider=minikube
					--k8sprovider=kubeadm
					`)

					break
				}

				// For Kubernetes 1.19 and above
				cmdStr = fmt.Sprintf("kubectl create ns %s --dry-run=client -o yaml | kubectl apply -f -", common.CMK8sNamespace)
				// For Kubernetes 1.18 and below
				//cmdStr = fmt.Sprintf("kubectl create ns %s --dry-run -o yaml | kubectl apply -f -", common.CMK8sNamespace)
				common.SysCall(cmdStr)

				// cmdStr = fmt.Sprintf("helm install --namespace %s %s -f %s ../helm-chart --debug", common.CMK8sNamespace, common.CMHelmReleaseName, common.FileStr)
				// if strings.ToLower(k8sprovider) == "gke" {
				// 	cmdStr += " --set metricServer.enabled=false"
				// }
				// //fmt.Println(cmdStr)
				// common.SysCall(cmdStr)

				if strings.ToLower(root.K8sprovider) == "gke" || strings.ToLower(root.K8sprovider) == "eks" || strings.ToLower(root.K8sprovider) == "aks" {
					cmdStr = fmt.Sprintf("helm install --namespace %s %s -f %s ../helm-chart --debug", common.CMK8sNamespace, common.CMHelmReleaseName, common.FileStr)
					cmdStr += " --set cb-restapigw.service.type=LoadBalancer"
					cmdStr += " --set cb-webtool.service.type=LoadBalancer"

					if strings.ToLower(root.K8sprovider) == "gke" || strings.ToLower(root.K8sprovider) == "aks" {
						cmdStr += " --set metricServer.enabled=false"
					}

					common.SysCall(cmdStr)
				} else {
					cmdStr = fmt.Sprintf("helm install --namespace %s %s -f %s ../helm-chart --debug", common.CMK8sNamespace, common.CMHelmReleaseName, common.FileStr)
					common.SysCall(cmdStr)
				}
			default:

			}

		}

	},
}

func init() {
	k8sCmd.AddCommand(runCmd)

	pf := runCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
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
