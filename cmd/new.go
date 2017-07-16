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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/telecoda/pico-go/console"
	"github.com/telecoda/pico-go/generate"
)

func er(err interface{}) {
	msg := fmt.Sprintf("Error: %s\n", err)
	color.Red(msg)
	os.Exit(-1)
}

var projectName string
var consoleType string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new pico-go project",
	Long: `This command will create a new pico-go game project in the current directory. 

Run the command in an empty directory into your $GOPATH
Where you want the initial project code to be generated.append

eg. $GOPATH/src/github.com/<your-github-profile>/<project-name>

This will generate all the scaffolding files to run your project with pico-go.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pico-go Create new project:")
		if len(args) != 1 {
			er("`new` command needs a name for the project")
		}

		if err := generate.NewProject(args[0], consoleType); err != nil {
			fmt.Printf("Failed to create new pico project: %s\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVar(&consoleType, "type", console.PICO8, "default console type")

}
