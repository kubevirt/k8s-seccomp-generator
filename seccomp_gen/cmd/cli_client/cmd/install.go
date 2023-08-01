package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

  var InstallCmd = &cobra.Command{
  Use:   "install",
  Short: "Install the application on the kubernetes cluster",
  Long:  `Install the application`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("installing the application...")
  },
}
