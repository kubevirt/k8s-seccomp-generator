package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: Move this elsewhere
var SUPPORTED_DISTROS [1]string = [1]string{"centos-stream8"}

  
func NewInstallCommand() *cobra.Command {
  var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install the application on the kubernetes cluster",
    Long:  `Install the application`,
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("installing the application...")
    },
    Args: func(cmd *cobra.Command, args []string) error {
      // Exactly one arg should be present, not less and not more
      if len(args) != 1 {
        return fmt.Errorf("OS distribution (and only that) must be present as the argument.")  
      }
      // given arg should be valid
      for _,dist := range SUPPORTED_DISTROS {
        if dist == args[0]{
          return nil
        }
      }
      return fmt.Errorf("Given OS distribution '%s' is invalid (or) not yet supported.", args[0])
    },  
  }

  return installCmd
}

