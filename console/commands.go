package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// init command list
var commands = map[string]Command{
	"HELP":  NewHelpCommand(),
	"DIR":   NewDirCommand(),
	"LS":    NewDirCommand(),
	"CD":    NewCDCommand(),
	"MKDIR": NewMkDirCommand(),
	"RUN":   NewRunCommand(),
}

type command struct {
	Name string
	Desc string
	Help string
}

type Command interface {
	Exec(pb PixelBuffer, statement string) error
}

/*
 Implemented commands
*/

type helpCommand struct {
	command
}

type dirCommand struct {
	command
}

type cdCommand struct {
	command
}

type mkdirCommand struct {
	command
}

type runCommand struct {
	command
}

func NewHelpCommand() Command {
	c := &helpCommand{
		command{

			Name: "HELP",
			Desc: "Display help",
		},
	}
	return c
}

func (h *helpCommand) Exec(pb PixelBuffer, statement string) error {
	pb.Color(BLUE)
	pb.Print("COMMANDS")
	pb.Print("")
	pb.Color(LIGHT_GRAY)
	pb.Print("LOAD <FILENAME>  SAVE <FILENAME>")
	pb.Print("RUN              RESUME")
	pb.Print("SHUTDOWN         REBOOT")
	pb.Print("CD <DIRNAME>     MKDIR <DIRNAME>")
	pb.Print("CD ..      TO GO UP A DIRECTORY")
	pb.Print("KEYCONFIG  TO CHOOSE BUTTONS")
	pb.Color(PINK)
	pb.Print("SPLORE     TO EXPLORE CARTRIGDES")
	pb.Print("")
	pb.Color(LIGHT_GRAY)
	pb.Print("PRESS ESC TO TOGGLE EDITOR VIEW")
	pb.Print("ALT+ENTER TO TOGGLE FULLSCREEN")
	pb.Print("CTRL-Q TO FASTQUIT")
	pb.Color(BLUE)
	pb.Print("SEE PICOGO.TXT FOR MORE INFO")
	pb.Print("OR VISIT WWW.??????.COM")
	return nil
}

func NewDirCommand() Command {
	c := &dirCommand{
		command{

			Name: "DIR",
			Desc: "List directory",
		},
	}
	return c
}

func pwd(pb PixelBuffer) {
	pb.Color(BLUE)
	dirStr := fmt.Sprintf("DIRECTORY: %s", _console.currentDir)
	pb.Print(dirStr)
}

func (d *dirCommand) Exec(pb PixelBuffer, statement string) error {
	pwd(pb)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := strings.ToUpper(file.Name())
		if file.IsDir() {
			pb.Color(PINK)
			pb.Print(name)
		} else {
			pb.Color(BLUE)
			pb.Print(name)
		}

	}
	return nil
}

func NewCDCommand() Command {
	c := &cdCommand{
		command{

			Name: "CD",
			Desc: "Change directory",
		},
	}
	return c
}

func (c *cdCommand) Exec(pb PixelBuffer, statement string) error {

	// split statement into tokens
	tokens := strings.Split(statement, " ")
	if len(tokens) < 2 {
		// no change
		return nil
	} else {
		newDir := strings.ToLower(tokens[1])

		newDir = _console.currentDir + "/" + newDir + "/"
		// TODO get code working for case insensitive dir names
		if err := os.Chdir(_console.baseDir + newDir); err != nil {
			pb.Color(WHITE)
			pb.Print("DIRECTORY NOT FOUND")
			return nil
		}

		// get dir details from full dir path
		fullDir, err := os.Getwd()
		if err != nil {
			return err
		}

		if len(fullDir) < len(_console.baseDir) {
			// you can't cd .. above starting dir
			os.Chdir(_console.baseDir + _console.currentDir)
			pb.Color(WHITE)
			pb.Print("CD: FAILED")
			return nil

		}

		if fullDir == _console.baseDir {
			_console.currentDir = "/"
		} else {
			parts := strings.Split(fullDir, "/")
			lastDir := parts[len(parts)-1]
			_console.currentDir = "/" + lastDir + "/"
		}

	}
	pwd(pb)

	return nil
}

func NewMkDirCommand() Command {
	c := &mkdirCommand{
		command{

			Name: "MKDIR",
			Desc: "Make directory",
			Help: "MAKE [NAME]",
		},
	}
	return c
}

func (m *mkdirCommand) Exec(pb PixelBuffer, statement string) error {

	// split statement into tokens
	tokens := strings.Split(statement, " ")
	if len(tokens) < 2 {
		pb.Color(LIGHT_GRAY)
		pb.Print(m.Help)
		return nil
	} else {
		newDir := strings.ToLower(tokens[1])

		newDir = _console.baseDir + _console.currentDir + newDir
		// TODO get code working for case insensitive dir names
		if err := os.Mkdir(newDir, 0700); err != nil {
			pb.Color(WHITE)
			pb.Print("MKDIR: FAILED")
			fmt.Printf("TEMP: err: %s\n", err)
			return nil
		}
	}

	return nil
}

func NewRunCommand() Command {
	c := &runCommand{
		command{

			Name: "RUN",
			Desc: "Run code",
		},
	}
	return c
}

func (r *runCommand) Exec(pb PixelBuffer, statement string) error {
	pb.Color(GREEN)
	pb.Print("RUNNING...")

	if runMode, ok := _console.modes[RUNTIME]; ok {
		runMode.Init()
		_console.SetMode(RUNTIME)
		return nil
	}
	panic("Unable to fetch RUNTIME mode!")

}
