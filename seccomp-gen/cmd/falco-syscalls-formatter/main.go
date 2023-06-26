package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
)

/*
  This binary takes the falco syscalls stream as the input and formats them and writes them to a file
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
  reader := bufio.NewReader(os.Stdin)
  syscallsMap := make(map[string]struct{})
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  for {
      select {
      case <-c:
        fmt.Println("Caught SIGINT, exiting...")
        for syscall := range syscallsMap {
          fmt.Println(syscall)
        }
        os.Exit(0)
      default:
        val := falcoOutput{}
        text, _ := reader.ReadBytes('\n')
        err := json.Unmarshal(text, &val)
        if err != nil {
          fmt.Println("Error : ", err)
        }
        syscall, ok :=  val.OutputFields["syscall.type"].(string)
        if ok && syscall != "" {
          syscallsMap[syscall] = struct{}{}
        }
      }
  }
}


