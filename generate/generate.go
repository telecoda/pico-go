package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/telecoda/pico-go/console"
)

func NewProject(projectName, consoleType string) error {
	fmt.Print("Creating new project\n")

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	consoleTypeStr, ok := console.ConsoleTypes[console.ConsoleType(consoleType)]
	if !ok {
		consoleTypes := make([]console.ConsoleType, len(console.ConsoleTypes))
		i := 0
		for t := range console.ConsoleTypes {
			consoleTypes[i] = t
			i++
		}
		return fmt.Errorf("Console type: %s not supported.  Options are %v", consoleType, consoleTypes)
	}

	projectPath := currDir + string(os.PathSeparator) + projectName

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Project will be created in directory: %s\n", projectPath)
	fmt.Printf("Enter y to continue\n")

	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if text[0] != 'y' && text[0] != 'Y' {
		return fmt.Errorf("You decided to cancel")
	}

	// calculate name of project to be generated
	goPath := build.Default.GOPATH

	if !strings.HasPrefix(projectPath, goPath) {
		return fmt.Errorf("Project directory must be within $GOPATH: %s", goPath)
	}

	// remove $GOPATH from projectPath
	repoPath := strings.TrimPrefix(projectPath, goPath+"/src/")

	// check directory does not exist
	files, err := ioutil.ReadDir(projectPath)
	if err == nil {
		// dir exists
		if len(files) > 0 {
			return fmt.Errorf("Directory %s already exists and contains files. Please select an empty directory", projectPath)
		}
	} else {
		// does not exist create new dir
		os.Mkdir(projectPath, 0777)
	}

	// create sub directories
	os.MkdirAll(projectPath+"/code", 0777)
	os.MkdirAll(projectPath+"/sprites", 0777)
	os.MkdirAll(projectPath+"/audio", 0777)

	filePath := projectPath + "/main.go"

	// generate main.go
	if err := generateMain(repoPath, filePath, consoleTypeStr); err != nil {
		return err
	}

	// duplicate template/code/cart.go
	fromSourcePath := goPath + "/src/github.com/telecoda/pico-go/template/code/cart.go"
	toSourcePath := projectPath + "/code/cart.go"
	demoSource, err := ioutil.ReadFile(fromSourcePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(toSourcePath, demoSource, 0666); err != nil {
		return err
	}
	// generate .gitignore
	gitIgnore := fmt.Sprintf("# ignore executables\n%s\n%s.exe\n# ignore state file\n.pico-go-state", projectName, string(os.PathSeparator))
	gitIgnorePath := projectPath + "/.gitignore"
	if err := ioutil.WriteFile(gitIgnorePath, []byte(gitIgnore), 0666); err != nil {
		return err
	}

	// copy sprites to local proj
	spritesFile := goPath + "/src/github.com/telecoda/pico-go/template/sprites/sprites.png"
	spriteBytes, err := ioutil.ReadFile(spritesFile)
	if err != nil {
		return err
	}
	localSpriteFile := fmt.Sprintf("./%s/sprites/sprites.png", projectName)
	err = ioutil.WriteFile(localSpriteFile, spriteBytes, 0600)
	if err != nil {
		return err
	}

	printBanner()
	// print statement to run code
	fmt.Printf("Congratulations you have created your pico-go project\ncd %s\npico-go run\n", projectName)
	return nil
}

func generateMain(projectPath, filePath, consoleType string) error {

	var data map[string]interface{}
	data = make(map[string]interface{})

	data["projectPath"] = projectPath
	data["consoleType"] = consoleType

	tmpl := template.New("")
	tmpl, err := tmpl.Parse(mainCodeTmpl)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)

	return ioutil.WriteFile(filePath, buf.Bytes(), 0666)

}

var mainCodeTmpl = `package main

import (
	"flag"

	"github.com/telecoda/pico-go/console"
	"{{ .projectPath }}/code"
)

/*

	THIS IS GENERATED CODE - DO NOT AMEND

*/

func main() {

	flag.Parse()

	// Create virtual console - based on cart config
	con, err := console.NewConsole(console.{{ .consoleType }})
	if err != nil {
		panic(err)
	}
	defer con.Destroy()

	cart := code.NewCart()

	if err := con.LoadCart(cart); err != nil {
		panic(err)
	}

	if err := con.Run(); err != nil {
		panic(err)
	}
}
`
