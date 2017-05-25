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
	"runtime"

	"time"

	"os"

	"path/filepath"

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

var lastBuilt time.Time
var currentCommand *exec.Cmd
var prevCommand *exec.Cmd
var restartChan chan bool
var quitChan chan bool
var cartBinary string

func run() {
	restartChan = make(chan bool, 1)
	quitChan = make(chan bool, 1)

	// binary name is based on package dir
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init binary name
	base := filepath.Base(dir)

	cartBinary = base
	if runtime.GOOS == "windows" {
		cartBinary += cartBinary + ".exe"
	} else {
		cartBinary = "./" + cartBinary
	}

	//  add a file watcher to main.go
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// initial code build
	if output, err := rebuildCode(); err != nil {
		if err != nil {
			fmt.Printf("Error trying to compile: %s\n", err)
			fmt.Printf("Code does not compiled: %s\n", output)
		}
		return
	}

	// if we're here code compiles so lets watch source file and run it

	// check for any code changes in code package
	err = watcher.Add("./code")
	if err != nil {
		log.Fatal(err)
	}

	// sourcecode watcher - checks for code changes
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					if output, err := rebuildCode(); err != nil {
						if err != nil {
							fmt.Printf("Error trying to compile: %s\n", err)
							fmt.Printf("Code does not compiled: %s\n", output)
						}
						continue
					}
					// code compiled - kill existing process
					restartChan <- true

				}
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	// Main loop
	// start 1st time
	restartChan <- true

	for {
		// sit and wait for stuff to happen..
		select {
		case <-quitChan:
			return
		case <-restartChan:
			restartBinary()

		}
	}
}

func restartBinary() error {
	if currentCommand != nil {
		prevCommand = currentCommand
		// start new binary
		if err := startBinary(); err != nil {
			return err
		}
		// then stop old one
		if err := stopBinary(prevCommand); err != nil {
			return err
		}

	} else {
		// this is first time
		if err := startBinary(); err != nil {
			return err
		}
	}

	return nil
}

func startBinary() error {
	// just start new binary
	currentCommand = exec.Command(cartBinary)
	currentCommand.Stderr = os.Stderr
	currentCommand.Stdout = os.Stdout
	// run compiled code - but don't wait
	err := currentCommand.Start()
	if err != nil {
		return err
	}
	// a little sleep to allow process to start
	time.Sleep(1 * time.Second)
	go func() {
		// wait till binary finishes
		currentCommand.Wait()
		if currentCommand != nil && currentCommand.ProcessState != nil && currentCommand.ProcessState.Success() {
			// clean exit, lets get out of here
			// probably a ctrl+q or closed window
			quitChan <- true
		}
	}()

	return nil
}

func stopBinary(command *exec.Cmd) error {
	return command.Process.Kill()
}

// rebuildCode will build a new exe
func rebuildCode() ([]byte, error) {

	// check how recently code was last built

	now := time.Now()

	diff := now.Sub(lastBuilt)

	if diff < time.Duration(5*time.Second) {
		// Just rebuild very recently.. ignoring
		return nil, nil
	}

	lastBuilt = now
	command := exec.Command("go", "build")
	return command.CombinedOutput()
}
