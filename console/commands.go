package console

// init command list
var commands = map[string]Command{
	"HELP": NewHelpCommand(),
}

type command struct {
	Name string
	Desc string
}

type Command interface {
	Exec(pb PixelBuffer) error
}

/*
 Implemented commands
*/

type helpCommand struct {
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

func (h *helpCommand) Exec(pb PixelBuffer) error {
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
	pb.Print("ALT+F4 OR CTRL-Q TO FASTQUIT")
	pb.Color(BLUE)
	pb.Print("SEE PICOGO.TXT FOR MORE INFO")
	pb.Print("OR VISIT WWW.??????.COM")
	return nil
}
