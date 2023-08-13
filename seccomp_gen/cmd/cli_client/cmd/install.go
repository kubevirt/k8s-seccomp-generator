package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	install "github.com/sudo-NithishKarthik/syscalls-tracer/pkg/install"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewInstallCommand() *cobra.Command {
	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install the application on the kubernetes cluster",
		Long:  `Install the application`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("installing the application...")
			distro := install.DistroFromString(args[0])
			fmt.Println("Selected Distro: ", distro.String())
			// see if the kubeconfig path is given in the flag
			kubeconfig, err := cmd.Flags().GetString("kube-config")
			if err != nil {
				panic(err.Error())
			}
			// get a kubernetes client
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err.Error())
			}
			// create the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				panic(err.Error())
			}
			// deploy and wait for the loader job to complete
			err = install.ConfigureNodes(clientset, distro)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Nodes have been configured.")
			// install tracer component manifests
			err = install.DeployTracerComponents(clientset)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Successfully deployed tracer components.")
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// Exactly one arg should be present, not less and not more
			if len(args) != 1 {
				return fmt.Errorf("OS distribution (and only that) must be present as the argument.")
			}
			// given arg should be valid
			for _, dist := range install.SUPPORTED_DISTROS {
				if dist == install.DistroFromString(args[0]) {
					return nil
				}
			}
			// TODO: Show the list of supported distros as well
			return fmt.Errorf("Given OS distribution '%s' is invalid (or) not yet supported.", args[0])
		},
	}

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	installCmd.Flags().String("kube-config", kubeconfig, "kuernetes config file")

	return installCmd
}
