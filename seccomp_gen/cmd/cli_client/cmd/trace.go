package cmd

import "github.com/spf13/cobra"

var traceCommand = &cobra.Command{
		Use:   "trace",
		Short: "Trace the syscalls made by kubernetes pods",
		Long:  `Trace the syscalls made by kubernetes pods`,
	}

func NewTraceAddCommand() *cobra.Command {
var addCommand = &cobra.Command{
		Use:   "add",
		Short: "Configure a trace to the database",
		Long:  `secgec-cli trace add $name $selector.\n $name is the name of the trace and $selector will be the one that determines which pods to trace`,
		Run: func(cmd *cobra.Command, args []string) {
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	
  return addCommand
}

func NewTraceStartCommand() *cobra.Command {
var addCommand = &cobra.Command{
		Use:   "start",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	
  return addCommand
}

func NewTraceCommand() *cobra.Command {
  return traceCommand
}

func init(){
  traceCommand.AddCommand(NewTraceAddCommand(), NewTraceStartCommand())
}
