package generate

import "fmt"
import "os"
import "bufio"

func NewProject(projectName string) error {
	fmt.Print("Creating new project\n")

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	projDir := currDir + string(os.PathSeparator) + projectName

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Project will be created in directory: %s\n", projDir)
	fmt.Printf("Enter y to continue\n")

	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if text != "y" && text != "Y" {
		return fmt.Errorf("You decided to cancel")
	}

	// check directory is empty

	return nil
}
