package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

/*
 * This binary takes the falco syscalls stream as the input (from stdio) and formats them and writes them to a file
 * Requirements: 
 * 1. Falco `json_output` has to be enabled
 * 2. Falco rule should have `syscall.type` in the output
*/

type falcoOutput struct {
  Hostname string `json:"hostname"`
  Output string `json:"output"`
  Priority string `json:"priority"`
  Rule string `json:"rule"`
  Source string `json:"source"`
  Time string `json:"time"`
  OutputFields map[string]interface{} `json:"output_fields"`
}

func main() {
  fmt.Println("Starting formatter...")
  scanner := bufio.NewScanner(os.Stdin)
  syscallsMap := make(map[string]struct{})
  c := make(chan os.Signal, 1)
  // we use SIGURG since Falco will be spawning this binary and it sends SIGURG when killed 
  // https://falco.org/docs/alerts/channels/#program-output - look at `Controlling the program output` section for more info
  signal.Notify(c, syscall.SIGURG)
  go func(c chan os.Signal){
    <-c
    fmt.Printf("Received SIGURG, exiting... ")
    writeSyscallsData(syscallsMap) 
    os.Exit(0)
  }(c)
  for scanner.Scan() {
        val := falcoOutput{}
        text := scanner.Text()
        err := json.Unmarshal([]byte(text), &val)
        if err != nil {
          fmt.Println("Cannot unmarshal input: ", err)
        }
        syscall, ok :=  val.OutputFields["syscall.type"].(string)
        if ok && syscall != "" {
        fmt.Println("Received syscall: ", syscall)
          syscallsMap[syscall] = struct{}{}
        }
  }
  if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

// write the syscalls to a data.json file
func writeSyscallsData(syscallsMap map[string]struct{}) {
    syscallData := make([]string, 0)
    for syscall := range syscallsMap {
      syscallData = append(syscallData, syscall)
    }
    file, _ := json.Marshal(syscallData)
    _ = os.WriteFile("/falco/data.json", file, 0644)
}
