package cmd

import "github.com/spf13/cobra"

var traceCommand = &cobra.Command{
		Use:   "trace",
		Short: "Trace the syscalls made by kubernetes pods",
		Long:  `Trace the syscalls made by kubernetes pods`,
	}

// We cannot store the map in-memory and hence we have to persist it in a file.
// This seems a bit too much for now. So we don't implement it now. 
func NewTraceAddCommand() *cobra.Command {
var addCommand = &cobra.Command{
		Use:   "configure",
		Short: "Add/update to a list of entities(called 'trace' list) that can be traced.",
		Long:  `secgen-cli trace add $name $selector.$name is the name of the trace and $selector will be the one that determines which pods to trace`,
		Run: func(cmd *cobra.Command, args []string) {
    // add/update it to a map of name:selector
    // this can also be used to update an existing entry of the map
		},
		Args: func(cmd *cobra.Command, args []string) error {
      // here we have to verify the selectors
			return nil
		},
	}

	
  return addCommand
}

func NewTraceStartCommand() *cobra.Command {
var startCommand = &cobra.Command{
		Use:   "start",
		Short: "Starts tracing an entity identified by the $selector",
		Long:  `This tool stores a list of entities on which we can start
            the trace by doing 'secgen-cli trace start $selector' and it
            will start tracing the syscalls made by the entity denoted by the $selector`,
		Run: func(cmd *cobra.Command, args []string) {
    // here we need the ability to communicate with the SYSCALLS_TRACER pod.
    // We need to send a request to the pod.
    // we will have a pkg specifically for communicating with the SYSCALLS_TRACER pod 
    // and it will be a layer of abstraction between the client and the SYSCALLS_TRACER server.
    // For now, we can just reach out to `localhost:30001` and it will get us the server.
    // but we need to keep the implementation modularized so that we can easily change it 
    // when we move to a different and more permanent solution for communication with the server.

    // NOTE: There can only be one trace going on at a time, so we need to verify that before 
    // starting the trace.
     
    // 1. Use the $selector to generate tracing.TracingConfiguration
    // 2. Send a request to localhost:30001/start with the tracing configuration as the body
		},
		Args: func(cmd *cobra.Command, args []string) error {
      // we need to verify whether the $selector is valid
			return nil
		},
	}
	
  return startCommand
}

func NewTraceCommand() *cobra.Command {
  return traceCommand
}

func init(){
  traceCommand.AddCommand(NewTraceAddCommand(), NewTraceStartCommand())
}
