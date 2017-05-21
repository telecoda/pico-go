// Copyright © 2017 @telecoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pico-go project",
	Long:  `This command will run the pico-go project in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run() {

	//  add a file watcher to main.go
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// initial code build
	if output, err := buildCode(); err != nil {
		if err != nil {
			fmt.Printf("Error trying to compile: %s\n", err)
			fmt.Printf("Code does not compiled: %s\n", output)
		}
		return
	}

	// if we're here code compiles so lets watch source file and run it

	// check for code changes
	err = watcher.Add("./code/cart.go")
	if err != nil {
		log.Fatal(err)
	}

	var command *exec.Cmd

	// sourcecode watcher - checks for code changes
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					if output, err := buildCode(); err != nil {
						if err != nil {
							fmt.Printf("Error trying to compile: %s\n", err)
							fmt.Printf("Code does not compiled: %s\n", output)
						}
						continue
					}
					// killCurrentProcess()
					// code compiled - kill existing process
					fmt.Printf("Killing process: %d\n", command.Process.Pid)
					syscall.Kill(command.Process.Pid, syscall.SIGKILL)
					fmt.Printf("Process killed\n")
				}
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	// Main loop
	hasQuit := false
	for !hasQuit {
		// TODO handle filename correctly hardcoded for now..
		command = exec.Command("./sprites_ex")
		// run compiled code
		err = command.Run()
		if command.ProcessState.Success() {
			// clean exit, lets get out of here
			// probably a ctrl+q or closed window
			hasQuit = true
		}
	}

	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		return
	}
}

// buildCode will build a new exe
func buildCode() ([]byte, error) {
	command := exec.Command("go", "build")
	return command.CombinedOutput()
}
