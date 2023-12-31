package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sudo-NithishKarthik/syscalls-tracer/pkg/tracing"
	"github.com/sudo-NithishKarthik/syscalls-tracer/pkg/seccomp"
)


// TODO : Add support for multiple nodes
// For multiple nodes, we need to first find a way for communication between the syscalls-tracer pod and the cli client

const (
  SYSCALLS_TRACER_POD_PORT = "30001"
)

var traceCommand = &cobra.Command{
	Use:   "trace",
	Short: "Trace the syscalls made by kubernetes pods",
	Long:  `Trace the syscalls made by kubernetes pods`,
}


// TO BE IMPLEMENTED LATER
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

			return nil
		},
	}
	return addCommand
}

func NewTraceStopCommand() *cobra.Command {
    var stopCommand = &cobra.Command{
		Use:   "stop",
		Short: "Stops tracing ",
		Long: `Stops tracing and sends the list of syscalls made by the pod`,
		Run: func(cmd *cobra.Command, args []string) {
      // 1, Send stop request
      resp, err := http.Get("http://localhost:"+SYSCALLS_TRACER_POD_PORT+"/stop")
			if err != nil {
				fmt.Println("Cannot send the stop request: ", err)
				return
      }
      // 2. Get data.json from the pod
      resp, err = http.Get("http://localhost:"+SYSCALLS_TRACER_POD_PORT+"/data.json")
			if err != nil {
				fmt.Println("Cannot send the stop request: ", err)
				return
      }
      syscallsData, err := io.ReadAll(resp.Body)
      resp.Body.Close()
			if err != nil {
				fmt.Println("Cannot read the response body: ", err)
				return
      }
      var syscalls []string
      err = json.Unmarshal(syscallsData, &syscalls)
			if err != nil {
				fmt.Println("Cannot unmarshal syscalls data: ", err)
				return
      }
      // 3. Generate seccomp profile from the syscalls
      profile := seccomp.SeccompProfile{}
      profile.DefaultAction = seccomp.ActErrno
      archs := make([]seccomp.Arch, 0)
      archs = append(archs, seccomp.ArchX86_64)
      archs = append(archs, seccomp.ArchX86)
      archs = append(archs, seccomp.ArchX32)
      profile.Architectures = archs
      
      scalls := make([]seccomp.Syscall, 0)
      for _, syscall := range syscalls {
        s := seccomp.Syscall{}
        s.Name = syscall
        s.Action = seccomp.ActAllow
        scalls = append(scalls, s)
      }
      profile.Syscalls = scalls 

      fmt.Printf("Profile: %v", profile)
    },
}
  return stopCommand
}
 
func NewTraceStartCommand() *cobra.Command {
	var startCommand = &cobra.Command{
		Use:   "start",
		Short: "Starts tracing an entity identified by the $selector",
		Long: `This tool stores a list of entities on which we can start
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
			selector := args[0]
			tracingConf := tracing.TracingConfiguration{}
			if strings.Contains(selector, "pod.name=") {
				parts := strings.Split(selector, "=")
				tracingConf.PodName = parts[1]
			}
			if strings.Contains(selector, "container.name=") {
				parts := strings.Split(selector, "=")
				tracingConf.ContainerName = parts[1]
			}
			if strings.Contains(selector, "pod.label.") {
				labelValuePair := selector[10:]
				parts := strings.Split(labelValuePair, "=")
				tracingConf.PodLabel = map[string]string{
					parts[0]: parts[1],
				}
			}
			// 2. Send a request to localhost:30001/start with the tracing configuration as the body
			jsonBody, err := json.Marshal(tracingConf)
			if err != nil {
				fmt.Println("Cannot marshal the tracing configuration: ", err)
				return
			}
			request, err := http.NewRequest("POST", "http://localhost:"+ SYSCALLS_TRACER_POD_PORT +"/start", bytes.NewBuffer(jsonBody))
			if err != nil {
				fmt.Println("Cannot form the request: ", err)
				return
			}
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				fmt.Println("Cannot send the request: ", err)
				return
			}
			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)
			fmt.Println("Started tracing...", string(body))
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// we need to verify whether the $selector is valid
			// Exactly one arg should be present, not less and not more
			if len(args) != 1 || args[0] == "" {
				return fmt.Errorf("selector must be present as the argument.")
			}
			// Supported selector types:
			// pod.name=$name
			// container.name=$name
			// pod.label.$label=$value
			exp, err := regexp.Compile("(^pod.name=)|(^pod.label.)|(^container.name=)")
			if err != nil {
				return fmt.Errorf("Regex error: %s", err)
			}
			res := exp.FindString(args[0])
			if res == "" {
				return fmt.Errorf("The selector provided: %s is invalid.", args[0])
			}
			return nil
		},
	}

	return startCommand
}

func NewTraceCommand() *cobra.Command {
	return traceCommand
}

func init() {
	traceCommand.AddCommand(NewTraceAddCommand(), NewTraceStartCommand(), NewTraceStopCommand())
}
