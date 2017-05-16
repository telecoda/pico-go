package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func NewProject(projectName string) error {
	fmt.Print("Creating new project\n")

	currDir, err := os.Getwd()
	if err != nil {
		return err
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
	goPath, _ := os.LookupEnv("GOPATH")

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
	os.MkdirAll(projectPath+"/fonts", 0777)
	os.MkdirAll(projectPath+"/sprites", 0777)
	os.MkdirAll(projectPath+"/audio", 0777)

	filePath := projectPath + "/main.go"

	// generate main.go
	if err := generateMain(repoPath, filePath); err != nil {
		return err
	}

	// duplicate demo/code/cart.go
	fromSourcePath := goPath + "/src/github.com/telecoda/pico-go/demo/code/cart.go"
	toSourcePath := projectPath + "/code/cart.go"
	demoSource, err := ioutil.ReadFile(fromSourcePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(toSourcePath, demoSource, 0666); err != nil {
		return err
	}
	// generate .gitignore

	gitIgnore := fmt.Sprintf("# ignore executables\n%s\n%s.exe\n", projectName, projectName)
	gitIgnorePath := projectPath + "/.gitignore"
	if err := ioutil.WriteFile(gitIgnorePath, []byte(gitIgnore), 0666); err != nil {
		return err
	}
	printBanner()
	// print statement to run code
	fmt.Printf("Congratulations you have created your pico-go project\ncd %s\ngo run main.go\n", projectName)
	return nil
}

func generateMain(projectPath, filePath string) error {

	var data map[string]interface{}
	data = make(map[string]interface{})

	data["projectPath"] = projectPath

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

func main() {

	flag.Parse()

	cart := code.NewCart()

	// Create virtual console - based on cart config
	con, err := console.NewConsole(cart.GetConfig())
	if err != nil {
		panic(err)
	}

	defer con.Destroy()

	if err := con.LoadCart(cart); err != nil {
		panic(err)
	}

	if err := con.Run(); err != nil {
		panic(err)
	}
}`